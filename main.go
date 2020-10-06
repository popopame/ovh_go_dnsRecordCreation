package main 

import (
	"os"
	"strconv"
	"github.com/ovh/go-ovh/ovh"
	"fmt"
)

//ovhzoneRecord define the struct that contain all the argument passed to the post api call to create a DNS record
type ovhZoneRecord struct {
	Id			int		`json:"id,omitempty"`
	FieldType	string	`json:"fieldType"`
	Subdomain	string	`json:"subDomain"`
	Target		string	`json:"target,omitempty"`
	TTL 		int 	`json:"ttl,omitempty"`
}


func main() {
	//Creation of the ovh client
	client, err := ovh.NewEndpointClient("ovh-eu")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}

	//Fetching all the needed info
	recordType, domain, subDomain, endpoint, actionRecord := setVariables()

	fmt.Printf("Going to %s the %s subdomain for %s domain pointing to %s \n",actionRecord,subDomain,domain,endpoint)

	//Switch case for handling the action
	switch actionRecord := os.Getenv("OVH_ACTION"); actionRecord {
	case "CREATE":
		record,err := createARecord(client,domain,recordType,subDomain,endpoint)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		}
		fmt.Printf("Success: %v\n", record)


	case "DELETE":
		err = deleteARecord(client,domain,recordType,endpoint)
		if err != nil {
			fmt.Printf("Error: %s",err)
		}
	
	}

}

//setVariable will fetch the environment variables containing the needed value for the script, if it fail to fetch one variable, it return an error and exit
func setVariables() (string, string, string, string, string){
	recordType := os.Getenv("OVH_RECORD_TYPE") //we test if the type of the record is set, if not, we fallback to A
	if len(recordType) == 0 {
		recordType = "A"
		fmt.Println("variables OVH_RECORD_TYPE not set, fallbacking to A (IPV4 type)")
	}


	domain := os.Getenv("OVH_DOMAIN")
	if len(domain) == 0 {
		fmt.Fprintf(os.Stderr,"Variable OVH_DOMAIN NOT SET, Please verify your environment variables")
		os.Exit(1)
	}
	subDomain := os.Getenv("OVH_SUBDOMAIN")
	if len(subDomain) == 0 {
		fmt.Fprintf(os.Stderr,"Variable OVH_SUBDOMAIN NOT SET, Please verify your environment variables")
		os.Exit(1)
	}
	endpoint := os.Getenv("OVH_IP_ENDPOINT")
	if len(endpoint) == 0 {
		fmt.Fprintf(os.Stderr,"Variable OVH_IP_ENDPOINT NOT SET, Please verify your environment variables")
		os.Exit(1)
	}
	actionRecord := os.Getenv("OVH_ACTION")
	if len(actionRecord) == 0 {
		fmt.Fprintf(os.Stderr,"Variable OVH_ACTION NOT SET, Please verify your environment variables")
		os.Exit(1)
	}

	return recordType,domain,subDomain,endpoint,actionRecord
}

//createARecord take as argument the client , zonename,diledtype,subdomain and target, and post a GET request on the url, and write the response to the record struc
func createARecord(ovhClient *ovh.Client, zoneName, fieldType, subdomain, target string) (*ovhZoneRecord, error) {
	url := "/domain/zone/"+ zoneName +"/record"

	parameters := ovhZoneRecord{
		FieldType: fieldType,
		Subdomain: subdomain,
		Target:		target,
	}

	record := ovhZoneRecord{}

	err := ovhClient.Post(url, &parameters, &record)
	if err != nil {
		return nil, fmt.Errorf("OVH API Call Failed: POST %s - %v \n with param %v", url, err, parameters)
	}

	return &record,nil
}

//getRecordID take as arg client, zonename, filedtype and subdomain, and query the id for the subdomain, the id is needed for the deletion of the record 
func getRecordId(ovhClient *ovh.Client, zoneName, fieldType, subdomain string) ([]int,error) {
	url := "/domain/zone/"+ zoneName + "/record?fieldType=" + fieldType + "&subDomain=" + subdomain
	ids := []int{}

	err := ovhClient.Get(url, &ids)
	if err != nil {
		return nil , fmt.Errorf("OVH API Call Failed: GET %s \n Error: %v", url, err)
	}

	return ids, err
}


//deleteARecord takes as argument the client, zoneName, fieldType and subdomain, and susing these arg it will first querry the id using the getRecordId func and make a post request with the id
func deleteARecord(ovhClient *ovh.Client, zoneName, fieldType, subdomain string) error {
	
	ids , _ := getRecordId(ovhClient, zoneName, fieldType, subdomain)
	for _, id := range ids {
		url := "/domain/zone/" + zoneName + "/record/" + strconv.Itoa(id)

		err := ovhClient.Delete(url, nil)
		if err != nil {
			return fmt.Errorf("OVH API Call Failed: DELETE %s \nError: %v",url,err)
		}
	}
	return nil
}
