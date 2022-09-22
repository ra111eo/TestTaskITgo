package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/ra111eo/TestTaskITgo/pb"
)

const admin = "admin"
const password = "password"
const dbname = "walletdb"
const localhost = "127.0.0.1"

type Ids struct {
	Id string `json:"id"`
}

type DBInfo struct {
	Total_rows int   `json:"total_rows"`
	Rows       []Ids `json:"rows"`
}

func makeDBurl() string {
	return fmt.Sprintf("http://%s:%s@%s:5984/%s", admin, password, localhost, dbname)
}

func GetDBinfo() (dbInfo *DBInfo, err error) {
	dbUrl := fmt.Sprintf("%s/_all_docs", makeDBurl())
	resp, err := http.Get(dbUrl)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	err = json.Unmarshal(body, &dbInfo)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	return dbInfo, nil
}

func UpdateEwalletFromDB(ewallet *pb.Ewallet) error {
	dbUrl := makeDBurl()
	jsonData, err := getInfoFromDB(dbUrl, ewallet.Id)
	if err != nil {
		log.Fatalln(err)
		return err
	}
	err = json.Unmarshal(jsonData, &ewallet)
	if err != nil {
		log.Fatalln(err)
		return err
	}
	return nil
}

func getInfoFromDB(dbUrl string, id string) ([]byte, error) {
	urlData := fmt.Sprintf("%s/%s", dbUrl, id)
	resp, err := http.Get(urlData)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	return body, nil
}

// used that except uuid.New
func GetNewId() (string, error) {
	url := fmt.Sprintf("http://%s:%s@%s:5984/_uuids", admin, password, localhost)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
		return "", err
	}
	fmt.Println(string(body))
	type Uuid struct {
		Uuids []string
	}
	var newId Uuid
	err = json.Unmarshal(body, &newId)
	if err != nil {
		log.Fatalln(err)
		return "", err
	}
	return newId.Uuids[0], nil
}

func UpdateDataInDB(ewallet *pb.Ewallet) error {
	dbUrl := makeDBurl()
	log.Printf("time to update data in db with id %s", ewallet.Id)
	ewalletURL := fmt.Sprintf("%s/%s", dbUrl, ewallet.Id)
	postBody, err := json.Marshal(ewallet)
	if err != nil {
		log.Fatalln(err)
		return err
	}
	req, err := http.NewRequest(http.MethodPut, ewalletURL, bytes.NewBuffer(postBody))
	if err != nil {
		log.Fatalln(err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	var client http.Client
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
		return err
	}
	sb := string(body)
	log.Printf(sb)
	return err
}
