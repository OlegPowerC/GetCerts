package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type Lmod struct {
	Lmod string `json:"lmod"`
}

type params struct {
	CertFile  string `json:"CertFile"`
	KeyFile   string `json:"KeyFile"`
	SrvUrl    string `json:"Url"`
	FullChain string `json:"FullChain"`
}

func GetCertKey(Dparam *params, Cfile string, fullpathname string, Outtoconsole bool) int {
	//Download key file
	var tr *http.Transport
	tr = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	clientkey := &http.Client{Transport: tr}

	reqprkey, _ := http.NewRequest("GET", Dparam.SrvUrl+"/"+Cfile, nil)
	respprkey, errprkey := clientkey.Do(reqprkey)
	if errprkey != nil {
		fmt.Println(errprkey)
	}

	dataprkey, errprkey := ioutil.ReadAll(respprkey.Body)
	if Outtoconsole {
		fmt.Println(string(dataprkey))
	} else {
		fmt.Println("************** some data *****************")
	}

	cfile, err := os.Create(fullpathname)
	if err != nil {
		panic(err)
	}
	defer cfile.Close()
	_, err = cfile.Write(dataprkey)
	if err != nil {
		fmt.Println(err)
		return 1
	}
	return 0
}

var JParams params

func main() {
	var lvmod Lmod

	// Открываем файл с настройками
	jSettingsFile, err := os.Open("settings.json")
	// Проверяем на ошибки
	if err != nil {
		fmt.Println("Error:", err)
	}
	defer jSettingsFile.Close()

	FData, err := ioutil.ReadAll(jSettingsFile)
	if err != nil {
		fmt.Println("Error:", err)
	}
	json.Unmarshal(FData, &JParams)

	//Открываем файл, содержащий дату последнего изменения файла на сервере
	_, errLmFileStatus := os.Stat("lastmodified.json")
	if errLmFileStatus != nil {
		if os.IsNotExist(err) {
			//Файл сертификата не найден
			print("File lastmodified.json not found but will be create")
		}
	} else {
		JSlfile, err := os.Open("lastmodified.json")
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			JsLdata, err := ioutil.ReadAll(JSlfile)
			if err != nil {
				fmt.Println("Error:", err)
			}

			//unmarshal to struct
			json.Unmarshal(JsLdata, &lvmod)
		}
		defer JSlfile.Close()
	}

	//cert filename
	_, err = os.Stat(JParams.CertFile)

	//Пока выставляем флаг в fase
	flagdownloadunc := false
	if os.IsNotExist(err) {
		//Файл сертификата не найден
		flagdownloadunc = true
	}
	fmt.Println("Download flag is:", flagdownloadunc)

	var tr *http.Transport
	tr = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	reqpr, _ := http.NewRequest("GET", JParams.SrvUrl+"/"+JParams.CertFile, nil)
	resppr, errpr := client.Do(reqpr)
	if errpr != nil {
		fmt.Println(errpr)
		panic("Exit")
	}

	//Смотрим когда файл последний раз модифицировался
	llastmodifiedhc := resppr.Header.Get("Last-Modified")
	llastmodifiedhc = strings.TrimSpace(llastmodifiedhc)

	//Смотрим когда мы последний раз обновляли локальный файл
	oldlmodjd := strings.TrimSpace(lvmod.Lmod)
	fmt.Println("Our date from JSON:", oldlmodjd, "date from header:", llastmodifiedhc)
	if llastmodifiedhc != oldlmodjd {
		fmt.Println("Date not equal!")
		JsonFile, err := os.Create("lastmodified.json")
		if err != nil {
			panic(err)
		}
		defer JsonFile.Close()

		lvmod.Lmod = llastmodifiedhc
		jst, _ := json.Marshal(&lvmod)
		fmt.Println(lvmod)
		fmt.Println(string(jst))
		JsonFile.Write(jst)
		flagdownloadunc = true
	}

	datapr, errpr := ioutil.ReadAll(resppr.Body)
	fmt.Println(string(datapr))
	if flagdownloadunc {
		//Нужно обновить файл сертификата но запрос заново можно не создавать, у нас уже есть содержимое в переменной: datapr
		fmt.Println("Need to renew cert")
		cfile, err := os.Create(JParams.CertFile)
		if err != nil {
			panic(err)
		}
		defer cfile.Close()
		_, err = cfile.Write(datapr)
		if err != nil {
			fmt.Println(err)
		}

	}
	if flagdownloadunc {
		fmt.Println("Need to renew fullchain and key")
		GetCertKey(&JParams, JParams.KeyFile, JParams.KeyFile, false)
		GetCertKey(&JParams, JParams.FullChain, JParams.FullChain, true)
		//Выполняем шелл скрипт для перезапуска WEB сервера
		out, err := exec.Command("/bin/sh", "httpdrestart.sh").CombinedOutput()
		if err != nil {
			fmt.Println("Error after execute shell script:", err)
		}
		fmt.Println("BashOutput:", string(out))
	}
}
