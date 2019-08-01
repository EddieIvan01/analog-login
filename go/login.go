package cumt_login

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/eddieivan01/nic"
)

type Loginer struct {
	uname string
	pwd   string
	S     *nic.Session
	t     string
}

const url = "http://jwxt.cumt.edu.cn/jwglxt"

var csrftoken, _ = regexp.Compile(`name="csrftoken" value="([0-9a-z\-,]+)"\/>`)

func NewLoginer(uname, pwd string) Loginer {
	return Loginer{uname, pwd, &nic.Session{}, getTimeStamp()}
}

func getTimeStamp() string {
	t := time.Now().Unix()
	return strconv.FormatInt(t, 10)
}

func (l Loginer) getPubKey() (map[string]string, error) {
	resp, err := l.S.Get(url+"/xtgl/login_getPublicKey.html?time="+l.t, nil)
	if err != nil {
		return nil, err
	}

	pub := make(map[string]string)
	err = resp.JSON(&pub)
	return pub, nil
}

func (l Loginer) getCSRFToken() (string, error) {
	resp, err := l.S.Get(url+"/xtgl/login_slogin.html?language=zh_CN&_t="+l.t, nil)
	if err != nil {
		return "", err
	}

	token := csrftoken.FindStringSubmatch(resp.Text)
	return token[1], nil
}

func bytes2Int(b []byte) *big.Int {
	r := big.NewInt(0)
	length := len(b)

	for i := 0; i < length; i++ {
		if b[i] == 0 {
			continue
		}
		t := big.NewInt(int64(b[i]))

		qt := (length - i - 1) / 4
		rem := (length - i - 1) % 4
		for j := 0; j < qt; j++ {
			t = t.Mul(t, big.NewInt(int64(math.Pow(2, 32))))
		}
		t = t.Mul(t, big.NewInt(int64(math.Pow(2, 8*float64(rem)))))
		r = r.Add(r, t)
	}
	return r
}

func (l Loginer) rsaEncrypt(pubPair map[string]string) ([]byte, error) {
	e := pubPair["exponent"]
	n := pubPair["modulus"]

	exponent, _ := base64.StdEncoding.DecodeString(e)
	modulus, _ := base64.StdEncoding.DecodeString(n)

	exponent = append([]byte{0}, exponent...)
	ee := binary.BigEndian.Uint32(exponent)
	nn := bytes2Int(modulus)

	pub := rsa.PublicKey{nn, int(ee)}
	return rsa.EncryptPKCS1v15(rand.Reader, &pub, []byte(l.pwd))
}

func (l Loginer) Login() error {
	csrftoken, err := l.getCSRFToken()
	if err != nil {
		return err
	}

	pubPair, err := l.getPubKey()
	if err != nil {
		return err
	}

	encPwd, err := l.rsaEncrypt(pubPair)
	if err != nil {
		return err
	}
	pwd := base64.StdEncoding.EncodeToString(encPwd)

	resp, err := l.S.Post(url+"/xtgl/login_slogin.html", &nic.H{
		Data: nic.KV{
			"yhm":       l.uname,
			"mm":        pwd,
			"csrftoken": csrftoken,
		},
		AllowRedirect: true,
	})
	if err != nil {
		return err
	}

	if strings.Contains(resp.Text, "教学评价") {
		return nil
	}

	return fmt.Errorf("%s", "login error")
}

/*
func main() {
	obj := NewLoginer("user", "pass")
	err := obj.Login()
	if err != nil {
		fmt.Println(err.Error())
	}
}
*/
