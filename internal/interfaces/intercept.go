package interfaces

import (
	"bufio"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"proxy/internal/domain/entity"
	"strconv"
)

func (proxy *Proxy) Intercept(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodConnect {
		proxy.tunnel(w, r)
		return
	}

	proxy.proxy(w, r)
}

func (proxy *Proxy) proxy(w http.ResponseWriter, r *http.Request) {
	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		proxy.logger.Error(err)
	}

	req := entity.Req{Request: string(dump), URL: r.RequestURI, Headers: r.Header}
	if err = proxy.dm.Insert(req); err != nil {
		fmt.Println("proxy.dm.Insert(req)")
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	resp, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		fmt.Println("http.DefaultTransport.RoundTrip(r)")
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			proxy.logger.Error(err)
		}
	}()

	copyHeaders(resp.Header, w.Header())

	w.WriteHeader(resp.StatusCode)
	if _, err = io.Copy(w, resp.Body); err != nil {
		proxy.logger.Error(err)
	}
}

func (proxy *Proxy) fillInformation(r *http.Request) {
	requestedUrl, _ := url.Parse(r.RequestURI)
	if requestedUrl.Scheme == "" {
		proxy.inf.Scheme = r.URL.Host
	} else {
		proxy.inf.Scheme = requestedUrl.Scheme
	}

	proxy.inf.InterceptedHttpsRequest = r
}

func (proxy *Proxy) tunnel(w http.ResponseWriter, r *http.Request) {
	proxy.fillInformation(r)

	hijackedConn, err := proxy.hijackConnect(w)
	if err != nil {
		proxy.logger.Error(err)
		return
	}
	defer hijackedConn.Close()

	TCPClientConn, err := proxy.initializeTCPClient(hijackedConn)
	if err != nil {
		proxy.logger.Error(err)
		return
	}
	defer TCPClientConn.Close()

	TCPServerConn, err := tls.Dial("tcp", proxy.inf.InterceptedHttpsRequest.Host, proxy.inf.Config)
	if err != nil {
		fmt.Println("tls.Dial(\"tcp\", proxy.inf.InterceptedHttpsRequest.Host, proxy.inf.Config)")
		proxy.logger.Error(err)
		return
	}

	err = proxy.doHttpsRequest(TCPClientConn, TCPServerConn)
	if err != nil {
		fmt.Println("doHttpsRequest(TCPClientConn, TCPServerConn)")
		proxy.logger.Error(err)
		return
	}

	dumped, err := httputil.DumpRequest(proxy.inf.ForwardedHttpsRequest, true)
	fmt.Println(string(dumped))
	if err != nil {
		proxy.logger.Error(err)
		return
	}

	err = proxy.dm.Insert(entity.Req{
		Headers: proxy.inf.ForwardedHttpsRequest.Header,
		Request: string(dumped),
		URL: fmt.Sprintf("https://%s%s", proxy.inf.ForwardedHttpsRequest.Host, proxy.inf.ForwardedHttpsRequest.URL.Path),
	})
	if err != nil {
		proxy.logger.Error(err)
		return
	}
}

func (proxy *Proxy) hijackConnect(w http.ResponseWriter) (net.Conn, error) {
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return nil, errors.New("hijacker !ok")
	}

	hijackedConn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return nil, err
	}

	_, err = hijackedConn.Write([]byte("HTTP/1.1 200 Connection established\r\n\r\n"))
	if err != nil {
		hijackedConn.Close()
		return nil, err
	}

	return hijackedConn, nil
}

func (proxy *Proxy) initializeTCPClient(hijackedConn net.Conn) (*tls.Conn, error) {
	cert, err := proxy.generateCertificate()
	if err != nil {
		return nil, err
	}

	proxy.inf.Config = &tls.Config{
		Certificates: []tls.Certificate{cert},
		ServerName:   proxy.inf.Scheme,
	}

	TCPClientConn := tls.Server(hijackedConn, proxy.inf.Config)
	err = TCPClientConn.Handshake()
	if err != nil {
		TCPClientConn.Close()
		hijackedConn.Close()
		fmt.Println("from Handshake")
		return nil, err
	}

	clientReader := bufio.NewReader(TCPClientConn)
	proxy.inf.ForwardedHttpsRequest, err = http.ReadRequest(clientReader)
	if err != nil {
		fmt.Println("from http.ReadRequest(clientReader)")
		return nil, err
	}

	return TCPClientConn, nil
}

func (proxy *Proxy) generateCertificate() (tls.Certificate, error) {
	rootDir, err := os.Getwd()
	if err != nil {
		return tls.Certificate{}, err
	}

	cmdGenDir := rootDir + "/genCerts"
	certsDir := cmdGenDir + "/certs/"
	certFilename := certsDir + proxy.inf.Scheme + ".crt"

	_, errStat := os.Stat(certFilename)
	if os.IsNotExist(errStat) {
		genCommand := exec.Command(cmdGenDir+"/gen_cert.sh", proxy.inf.Scheme, strconv.Itoa(rand.Intn(1000)))

		out, err := genCommand.CombinedOutput()
		if err != nil {
			fmt.Println(string(out))
			return tls.Certificate{}, err
		}
	}

	cert, err := tls.LoadX509KeyPair(certFilename, cmdGenDir+"/cert.key")
	if err != nil {
		fmt.Println("from tls.LoadX509KeyPair(certFilename, rootDir+\"/cert.key\")")
		return tls.Certificate{}, err
	}

	return cert, nil
}

func (proxy *Proxy) doHttpsRequest(TCPClientConn *tls.Conn, TCPServerConn *tls.Conn) error {
	rawReq, err := httputil.DumpRequest(proxy.inf.ForwardedHttpsRequest, true)
	_, err = TCPServerConn.Write(rawReq)
	if err != nil {
		fmt.Println("TCPServerConn.Write(rawReq)")
		return err
	}

	serverReader := bufio.NewReader(TCPServerConn)
	TCPServerResponse, err := http.ReadResponse(serverReader, proxy.inf.ForwardedHttpsRequest)
	if err != nil {
		fmt.Println("http.ReadResponse(serverReader, proxy.inf.ForwardedHttpsRequest)")
		return err
	}

	decodedResponse, err := decodeResponse(TCPServerResponse)
	if err != nil {
		fmt.Println("decodeResponse(TCPServerResponse)")
		return err
	}
	_, err = TCPClientConn.Write(decodedResponse)
	if err != nil {
		fmt.Println("TCPClientConn.Write(decodedResponse)")
		return err
	}

	return nil
}
