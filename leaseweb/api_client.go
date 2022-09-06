package leaseweb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

var (
	leasewebAPIURL   string
	leasewebAPIToken string
	leasewebClient   *http.Client
)

// Server -
type Server struct {
	ID       string
	Contract struct {
		Reference string
	}
	NetworkInterfaces struct {
		Public struct {
			IP string
		}
		RemoteManagement struct {
			IP string
		}
	}
	Location struct {
		Site  string
		Suite string
		Rack  string
		Unit  string
	}
}

// IP -
type IP struct {
	IP            string
	ReverseLookup string
	NullRouted    bool
}

// DHCPLease -
type DHCPLease struct {
	Leases []struct {
		IP       string
		Bootfile string
	}
}

// GetBootfile -
func (l *DHCPLease) GetBootfile() string {
	if len(l.Leases) == 0 {
		return ""
	}
	return l.Leases[0].Bootfile
}

// PowerInfo -
type PowerInfo struct {
	IPMI struct {
		Status string
	}
	PDU struct {
		Status string
	}
}

// IsPoweredOn -
func (p *PowerInfo) IsPoweredOn() bool {
	return p.PDU.Status == "on" // TODO also take ipmi into account
}

// NetworkInterfaceInfo -
type NetworkInterfaceInfo struct {
	Status string
}

// IsOpened -
func (n *NetworkInterfaceInfo) IsOpened() bool {
	return n.Status == "OPEN"
}

// NotificationSetting -
type NotificationSetting struct {
	ID        string  `json:"id,omitempty"`
	Frequency string  `json:"frequency"`
	Threshold float64 `json:"threshold,string"`
	Unit      string  `json:"unit"`
}

// Credential -
type Credential struct {
	Type     string `json:"type"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// OperatingSystem -
type OperatingSystem struct {
	ID   string
	Name string
}

// ControlPanel -
type ControlPanel struct {
	ID   string
	Name string
}

// Payload -
type Payload map[string]interface{}

// Job -
type Job struct {
	UUID    string
	Status  string
	Payload Payload
}

func getServer(serverID string) (*Server, error) {
	request, err := http.NewRequest("GET", fmt.Sprintf("%s/bareMetals/v2/servers/%s", leasewebAPIURL, serverID), nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("X-Lsw-Auth", leasewebAPIToken)

	response, err := leasewebClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error getting server, api response %v", response.StatusCode)
	}

	var server Server
	err = json.NewDecoder(response.Body).Decode(&server)
	if err != nil {
		return nil, err
	}

	server.NetworkInterfaces.Public.IP = strings.SplitN(server.NetworkInterfaces.Public.IP, "/", 2)[0]
	server.NetworkInterfaces.RemoteManagement.IP = strings.SplitN(server.NetworkInterfaces.RemoteManagement.IP, "/", 2)[0]

	return &server, nil
}

func getServerIP(serverID string, ip string) (*IP, error) {
	request, err := http.NewRequest("GET", fmt.Sprintf("%s/bareMetals/v2/servers/%s/ips/%s", leasewebAPIURL, serverID, ip), nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("X-Lsw-Auth", leasewebAPIToken)

	response, err := leasewebClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error getting ip, api response %v", response.StatusCode)
	}

	var ipData IP
	err = json.NewDecoder(response.Body).Decode(&ipData)
	if err != nil {
		return nil, err
	}

	return &ipData, nil
}

func getServerLease(serverID string) (*DHCPLease, error) {
	request, err := http.NewRequest("GET", fmt.Sprintf("%s/bareMetals/v2/servers/%s/leases", leasewebAPIURL, serverID), nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("X-Lsw-Auth", leasewebAPIToken)

	response, err := leasewebClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error getting leases, api response %v", response.StatusCode)
	}

	var dhcpLease DHCPLease
	err = json.NewDecoder(response.Body).Decode(&dhcpLease)
	if err != nil {
		return nil, err
	}

	return &dhcpLease, nil
}

func getPowerInfo(serverID string) (*PowerInfo, error) {
	request, err := http.NewRequest("GET", fmt.Sprintf("%s/bareMetals/v2/servers/%s/powerInfo", leasewebAPIURL, serverID), nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("X-Lsw-Auth", leasewebAPIToken)

	response, err := leasewebClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error getting leases, api response %v", response.StatusCode)
	}

	var powerInfo PowerInfo
	err = json.NewDecoder(response.Body).Decode(&powerInfo)
	if err != nil {
		return nil, err
	}

	return &powerInfo, nil
}

func getNetworkInterfaceInfo(serverID string, networkType string) (*NetworkInterfaceInfo, error) {
	request, err := http.NewRequest("GET", fmt.Sprintf("%s/bareMetals/v2/servers/%s/networkInterfaces/%s", leasewebAPIURL, serverID, networkType), nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("X-Lsw-Auth", leasewebAPIToken)

	response, err := leasewebClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error getting network interface info, api response %v", response.StatusCode)
	}

	var networkInterfaceInfo NetworkInterfaceInfo
	err = json.NewDecoder(response.Body).Decode(&networkInterfaceInfo)
	if err != nil {
		return nil, err
	}

	return &networkInterfaceInfo, nil
}

func updateReference(serverID string, reference string) error {
	requestBody := new(bytes.Buffer)
	err := json.NewEncoder(requestBody).Encode(struct {
		Reference string `json:"reference"`
	}{
		Reference: reference,
	})
	if err != nil {
		return err
	}

	request, err := http.NewRequest("PUT", fmt.Sprintf("%s/bareMetals/v2/servers/%s", leasewebAPIURL, serverID), requestBody)
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Lsw-Auth", leasewebAPIToken)

	response, err := leasewebClient.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("error updating reference, api response %v", response.StatusCode)
	}

	return nil
}

func updateReverseLookup(serverID string, ip string, reverseLookup string) error {
	requestBody := new(bytes.Buffer)
	err := json.NewEncoder(requestBody).Encode(struct {
		ReverseLookup string `json:"reverseLookup"`
	}{
		ReverseLookup: reverseLookup,
	})
	if err != nil {
		return err
	}

	request, err := http.NewRequest("PUT", fmt.Sprintf("%s/bareMetals/v2/servers/%s/ips/%s", leasewebAPIURL, serverID, ip), requestBody)
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Lsw-Auth", leasewebAPIToken)

	response, err := leasewebClient.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("error updating reverse lookup, api response %v", response.StatusCode)
	}

	return nil
}

func powerOnServer(serverID string) error {
	request, err := http.NewRequest("POST", fmt.Sprintf("%s/bareMetals/v2/servers/%s/powerOn", leasewebAPIURL, serverID), nil)
	if err != nil {
		return err
	}
	request.Header.Set("X-Lsw-Auth", leasewebAPIToken)

	response, err := leasewebClient.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusAccepted {
		return fmt.Errorf("error powering on server, api response %v", response.StatusCode)
	}

	return nil
}

func powerOffServer(serverID string) error {
	request, err := http.NewRequest("POST", fmt.Sprintf("%s/bareMetals/v2/servers/%s/powerOff", leasewebAPIURL, serverID), nil)
	if err != nil {
		return err
	}
	request.Header.Set("X-Lsw-Auth", leasewebAPIToken)

	response, err := leasewebClient.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusAccepted {
		return fmt.Errorf("error powering off server, api response %v", response.StatusCode)
	}

	return nil
}

func addDHCPLease(serverID string, bootfile string) error {
	requestBody := new(bytes.Buffer)
	err := json.NewEncoder(requestBody).Encode(struct {
		Bootfile string `json:"bootfile"`
	}{
		Bootfile: bootfile,
	})
	if err != nil {
		return err
	}

	request, err := http.NewRequest("POST", fmt.Sprintf("%s/bareMetals/v2/servers/%s/leases", leasewebAPIURL, serverID), requestBody)
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Lsw-Auth", leasewebAPIToken)

	response, err := leasewebClient.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("error adding dhcp lease, api response %v", response.StatusCode)
	}

	return nil
}

func removeDHCPLease(serverID string) error {
	request, err := http.NewRequest("DELETE", fmt.Sprintf("%s/bareMetals/v2/servers/%s/leases", leasewebAPIURL, serverID), nil)
	if err != nil {
		return err
	}
	request.Header.Set("X-Lsw-Auth", leasewebAPIToken)

	response, err := leasewebClient.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("error removing dhcp lease, api response %v", response.StatusCode)
	}

	return nil
}

func openNetworkInterface(serverID string, networkType string) error {
	request, err := http.NewRequest("POST", fmt.Sprintf("%s/bareMetals/v2/servers/%s/networkInterfaces/%s/open", leasewebAPIURL, serverID, networkType), nil)
	if err != nil {
		return err
	}
	request.Header.Set("X-Lsw-Auth", leasewebAPIToken)

	response, err := leasewebClient.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("error opening network interface, api response %v", response.StatusCode)
	}

	return nil
}

func closeNetworkInterface(serverID string, networkType string) error {
	request, err := http.NewRequest("POST", fmt.Sprintf("%s/bareMetals/v2/servers/%s/networkInterfaces/%s/close", leasewebAPIURL, serverID, networkType), nil)
	if err != nil {
		return err
	}
	request.Header.Set("X-Lsw-Auth", leasewebAPIToken)

	response, err := leasewebClient.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("error closing network interface, api response %v", response.StatusCode)
	}

	return nil
}

func nullIP(serverID string, IP string) error {
	request, err := http.NewRequest("POST", fmt.Sprintf("%s/bareMetals/v2/servers/%s/ips/%s/null", leasewebAPIURL, serverID, IP), nil)
	if err != nil {
		return err
	}
	request.Header.Set("X-Lsw-Auth", leasewebAPIToken)

	response, err := leasewebClient.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusAccepted {
		return fmt.Errorf("error nulling ip of the server, api response %v", response.StatusCode)
	}

	return nil
}

func unnullIP(serverID string, IP string) error {
	request, err := http.NewRequest("POST", fmt.Sprintf("%s/bareMetals/v2/servers/%s/ips/%s/unnull", leasewebAPIURL, serverID, IP), nil)
	if err != nil {
		return err
	}
	request.Header.Set("X-Lsw-Auth", leasewebAPIToken)

	response, err := leasewebClient.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusAccepted {
		return fmt.Errorf("error unnulling server ip of the server, api response %v", response.StatusCode)
	}

	return nil
}

func createDedicatedServerNotificationSetting(serverID string, notificationType string, notificationSetting *NotificationSetting) (*NotificationSetting, error) {
	requestBody := new(bytes.Buffer)
	err := json.NewEncoder(requestBody).Encode(notificationSetting)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", fmt.Sprintf("%s/bareMetals/v2/servers/%s/notificationSettings/%s", leasewebAPIURL, serverID, notificationType), requestBody)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Lsw-Auth", leasewebAPIToken)

	response, err := leasewebClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("error creating server notification setting, api response %v", response.StatusCode)
	}

	var createdNotificationSetting NotificationSetting
	err = json.NewDecoder(response.Body).Decode(&createdNotificationSetting)
	if err != nil {
		return nil, err
	}

	return &createdNotificationSetting, nil
}

func getDedicatedServerNotificationSetting(serverID string, notificationType string, notificationSettingID string) (*NotificationSetting, error) {
	request, err := http.NewRequest("GET", fmt.Sprintf("%s/bareMetals/v2/servers/%s/notificationSettings/%s/%s", leasewebAPIURL, serverID, notificationType, notificationSettingID), nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("X-Lsw-Auth", leasewebAPIToken)

	response, err := leasewebClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error getting server notification setting, api response %v", response.StatusCode)
	}

	var notificationSetting NotificationSetting
	err = json.NewDecoder(response.Body).Decode(&notificationSetting)
	if err != nil {
		return nil, err
	}

	return &notificationSetting, nil
}

func updateDedicatedServerNotificationSetting(serverID string, notificationType string, notificationSettingID string, notificationSetting *NotificationSetting) (*NotificationSetting, error) {
	requestBody := new(bytes.Buffer)
	err := json.NewEncoder(requestBody).Encode(notificationSetting)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("PUT", fmt.Sprintf("%s/bareMetals/v2/servers/%s/notificationSettings/%s/%s", leasewebAPIURL, serverID, notificationType, notificationSettingID), requestBody)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Lsw-Auth", leasewebAPIToken)

	response, err := leasewebClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error updating server notification setting, api response %v", response.StatusCode)
	}

	var updatedNotificationSetting NotificationSetting
	err = json.NewDecoder(response.Body).Decode(&updatedNotificationSetting)
	if err != nil {
		return nil, err
	}

	return &updatedNotificationSetting, nil
}

func deleteDedicatedServerNotificationSetting(serverID string, notificationType string, notificationSettingID string) error {
	request, err := http.NewRequest("DELETE", fmt.Sprintf("%s/bareMetals/v2/servers/%s/notificationSettings/%s/%s", leasewebAPIURL, serverID, notificationType, notificationSettingID), nil)
	if err != nil {
		return err
	}
	request.Header.Set("X-Lsw-Auth", leasewebAPIToken)

	response, err := leasewebClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("error deleting server notification setting, api response %v", response.StatusCode)
	}

	return nil
}

func createDedicatedServerCredential(serverID string, credential *Credential) (*Credential, error) {
	requestBody := new(bytes.Buffer)
	err := json.NewEncoder(requestBody).Encode(credential)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", fmt.Sprintf("%s/bareMetals/v2/servers/%s/credentials", leasewebAPIURL, serverID), requestBody)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Lsw-Auth", leasewebAPIToken)

	response, err := leasewebClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error creating server credential, api response %v", response.StatusCode)
	}

	var createdCredential Credential
	err = json.NewDecoder(response.Body).Decode(&createdCredential)
	if err != nil {
		return nil, err
	}

	return &createdCredential, nil
}

func getDedicatedServerCredential(serverID string, credentialType string, username string) (*Credential, error) {
	request, err := http.NewRequest("GET", fmt.Sprintf("%s/bareMetals/v2/servers/%s/credentials/%s/%s", leasewebAPIURL, serverID, credentialType, username), nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("X-Lsw-Auth", leasewebAPIToken)

	response, err := leasewebClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error getting server credential, api response %v", response.StatusCode)
	}

	var credential Credential
	err = json.NewDecoder(response.Body).Decode(&credential)
	if err != nil {
		return nil, err
	}

	return &credential, nil
}

func updateDedicatedServerCredential(serverID string, credential *Credential) (*Credential, error) {
	requestBody := new(bytes.Buffer)
	err := json.NewEncoder(requestBody).Encode(struct {
		Password string `json:"password"`
	}{
		Password: credential.Password,
	})
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("PUT", fmt.Sprintf("%s/bareMetals/v2/servers/%s/credentials/%s/%s", leasewebAPIURL, serverID, credential.Type, credential.Username), requestBody)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Lsw-Auth", leasewebAPIToken)

	response, err := leasewebClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error updating server credential, api response %v", response.StatusCode)
	}

	var updatedCredential Credential
	err = json.NewDecoder(response.Body).Decode(&updatedCredential)
	if err != nil {
		return nil, err
	}

	return &updatedCredential, nil
}

func deleteDedicatedServerCredential(serverID string, credential *Credential) error {
	request, err := http.NewRequest("DELETE", fmt.Sprintf("%s/bareMetals/v2/servers/%s/credentials/%s/%s", leasewebAPIURL, serverID, credential.Type, credential.Username), nil)
	if err != nil {
		return err
	}
	request.Header.Set("X-Lsw-Auth", leasewebAPIToken)

	response, err := leasewebClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("error deleting server credential, api response %v", response.StatusCode)
	}

	return nil
}

func getOperatingSystems() ([]OperatingSystem, error) {
	request, err := http.NewRequest("GET", fmt.Sprintf("%s/bareMetals/v2/operatingSystems", leasewebAPIURL), nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("X-Lsw-Auth", leasewebAPIToken)

	response, err := leasewebClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error getting operating systems, api response %v", response.StatusCode)
	}

	var operatingSystems struct {
		OperatingSystems []OperatingSystem
	}

	err = json.NewDecoder(response.Body).Decode(&operatingSystems)
	if err != nil {
		return nil, err
	}

	// to be exact we would need to support pagination by checking the metadata and make multiple requests if needed
	// but with the default offset and limit values we already get the full list at the moment

	return operatingSystems.OperatingSystems, nil
}

func getControlPanels(operatingSystemID string) ([]ControlPanel, error) {
	request, err := http.NewRequest("GET", fmt.Sprintf("%s/bareMetals/v2/controlPanels", leasewebAPIURL), nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("X-Lsw-Auth", leasewebAPIToken)
	if operatingSystemID != "" {
		q := request.URL.Query()
		q.Add("operatingSystemId", operatingSystemID)
		request.URL.RawQuery = q.Encode()
	}

	response, err := leasewebClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error getting control panels, api response %v", response.StatusCode)
	}

	var controlPanels struct {
		ControlPanels []ControlPanel
	}

	err = json.NewDecoder(response.Body).Decode(&controlPanels)
	if err != nil {
		return nil, err
	}

	return controlPanels.ControlPanels, nil
}

func launchInstallationJob(serverID string, payload *Payload) (*Job, error) {
	requestBody := new(bytes.Buffer)
	err := json.NewEncoder(requestBody).Encode(payload)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", fmt.Sprintf("%s/bareMetals/v2/servers/%s/install", leasewebAPIURL, serverID), requestBody)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Lsw-Auth", leasewebAPIToken)

	response, err := leasewebClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("error launching installation job, api response %v", response.StatusCode)
	}

	var installationJob Job

	err = json.NewDecoder(response.Body).Decode(&installationJob)
	if err != nil {
		return nil, err
	}

	return &installationJob, nil
}

func getLatestInstallationJob(serverID string) (*Job, error) {
	request, err := http.NewRequest("GET", fmt.Sprintf("%s/bareMetals/v2/servers/%s/jobs", leasewebAPIURL, serverID), nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("X-Lsw-Auth", leasewebAPIToken)

	q := request.URL.Query()
	q.Add("type", "install")
	request.URL.RawQuery = q.Encode()

	response, err := leasewebClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error getting latest installation job, api response %v", response.StatusCode)
	}

	var jobs struct {
		Jobs []Job
	}

	err = json.NewDecoder(response.Body).Decode(&jobs)
	if err != nil {
		return nil, err
	}

	return &jobs.Jobs[0], nil
}

func getJob(serverID string, jobUUID string) (*Job, error) {
	request, err := http.NewRequest("GET", fmt.Sprintf("%s/bareMetals/v2/servers/%s/jobs/%s", leasewebAPIURL, serverID, jobUUID), nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("X-Lsw-Auth", leasewebAPIToken)

	response, err := leasewebClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error getting job status, api response %v", response.StatusCode)
	}

	var job Job

	err = json.NewDecoder(response.Body).Decode(&job)
	if err != nil {
		return nil, err
	}

	return &job, nil
}

func getServersBatch(offset int, limit int, site string) ([]Server, error) {
	request, err := http.NewRequest("GET", fmt.Sprintf("%s/bareMetals/v2/servers", leasewebAPIURL), nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("X-Lsw-Auth", leasewebAPIToken)

	if offset >= 0 {
		q := request.URL.Query()
		q.Add("offset", strconv.Itoa(offset))
		request.URL.RawQuery = q.Encode()
	}

	if limit >= 0 {
		q := request.URL.Query()
		q.Add("limit", strconv.Itoa(limit))
		request.URL.RawQuery = q.Encode()
	}

	if site != "" {
		q := request.URL.Query()
		q.Add("site", site)
		request.URL.RawQuery = q.Encode()
	}

	response, err := leasewebClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error getting server list, api response %v", response.StatusCode)
	}

	var serverList struct {
		Servers []Server
	}

	err = json.NewDecoder(response.Body).Decode(&serverList)
	if err != nil {
		return nil, err
	}

	return serverList.Servers, nil
}

func getAllServers(site string) ([]Server, error) {
	var allServers []Server
	offset := 0
	limit := 20

	serversBatch, err := getServersBatch(offset, limit, site)
	if err != nil {
		return nil, err
	}
	allServers = append(allServers, serversBatch...)

	for len(serversBatch) != 0 {
		offset += limit
		serversBatch, err = getServersBatch(offset, limit, site)
		if err != nil {
			return nil, err
		}
		allServers = append(allServers, serversBatch...)
	}

	return allServers, nil
}
