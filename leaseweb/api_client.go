package leaseweb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
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
	ID        string `json:"id,omitempty"`
	Frequency string `json:"frequency"`
	Threshold string `json:"threshold"`
	Unit      string `json:"unit"`
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

func getServerMainIP(serverID string, mainIP string) (*IP, error) {
	request, err := http.NewRequest("GET", fmt.Sprintf("%s/bareMetals/v2/servers/%s/ips/%s", leasewebAPIURL, serverID, mainIP), nil)
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
		return nil, fmt.Errorf("error getting main ip, api response %v", response.StatusCode)
	}

	var ip IP
	err = json.NewDecoder(response.Body).Decode(&ip)
	if err != nil {
		return nil, err
	}

	return &ip, nil
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

func updateReverseLookup(serverID string, mainIP string, reverseLookup string) error {
	requestBody := new(bytes.Buffer)
	err := json.NewEncoder(requestBody).Encode(struct {
		ReverseLookup string `json:"reverseLookup"`
	}{
		ReverseLookup: reverseLookup,
	})
	if err != nil {
		return err
	}

	request, err := http.NewRequest("PUT", fmt.Sprintf("%s/bareMetals/v2/servers/%s/ips/%s", leasewebAPIURL, serverID, mainIP), requestBody)
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

func createDedicatedServerNotificationSettingBandwidth(serverID string, notificationSetting *NotificationSetting) (*NotificationSetting, error) {
	requestBody := new(bytes.Buffer)
	err := json.NewEncoder(requestBody).Encode(notificationSetting)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", fmt.Sprintf("%s/bareMetals/v2/servers/%s/notificationSettings/bandwidth", leasewebAPIURL, serverID), requestBody)
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

func getDedicatedServerNotificationSettingBandwidth(serverID string, notificationSettingID string) (*NotificationSetting, error) {
	request, err := http.NewRequest("GET", fmt.Sprintf("%s/bareMetals/v2/servers/%s/notificationSettings/bandwidth/%s", leasewebAPIURL, serverID, notificationSettingID), nil)
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

func updateDedicatedServerNotificationSettingBandwidth(serverID string, notificationSettingID string, notificationSetting *NotificationSetting) (*NotificationSetting, error) {
	requestBody := new(bytes.Buffer)
	err := json.NewEncoder(requestBody).Encode(notificationSetting)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("PUT", fmt.Sprintf("%s/bareMetals/v2/servers/%s/notificationSettings/bandwidth/%s", leasewebAPIURL, serverID, notificationSettingID), requestBody)
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

func deleteDedicatedServerNotificationSettingBandwidth(serverID string, notificationSettingID string) error {
	request, err := http.NewRequest("DELETE", fmt.Sprintf("%s/bareMetals/v2/servers/%s/notificationSettings/bandwidth/%s", leasewebAPIURL, serverID, notificationSettingID), nil)
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
