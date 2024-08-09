package to_domain_entity

import (
  "fmt"

  "github.com/leaseweb/leaseweb-go-sdk/publicCloud"
  "github.com/leaseweb/terraform-provider-leaseweb/internal/core/domain/public_cloud"
  "github.com/leaseweb/terraform-provider-leaseweb/internal/core/shared/enum"
  "github.com/leaseweb/terraform-provider-leaseweb/internal/core/shared/value_object"
  "github.com/leaseweb/terraform-provider-leaseweb/internal/repositories/shared"
)

// AdaptInstance adapts an instance domain entity to an sdk instance model.
func AdaptInstance(
  sdkInstance publicCloud.Instance,
) (
  *public_cloud.Instance,
  error,
) {
  var autoScalingGroup *public_cloud.AutoScalingGroup

  instanceId, err := value_object.NewUuid(sdkInstance.GetId())
  if err != nil {
    return nil, fmt.Errorf("AdaptInstance: %w", err)
  }

  state, err := enum.NewState(string(sdkInstance.GetState()))
  if err != nil {
    return nil, fmt.Errorf("AdaptInstance: %w", err)
  }

  rootDiskSize, err := value_object.NewRootDiskSize(
    int(sdkInstance.GetRootDiskSize()),
  )
  if err != nil {
    return nil, fmt.Errorf("AdaptInstance: %w", err)
  }

  rootDiskStorageType, err := enum.NewRootDiskStorageType(
    string(sdkInstance.GetRootDiskStorageType()),
  )
  if err != nil {
    return nil, fmt.Errorf("AdaptInstance: %w", err)
  }

  ips, err := adaptIps(sdkInstance.GetIps())
  if err != nil {
    return nil, fmt.Errorf("AdaptInstance:  %w", err)
  }

  contract, err := adaptContract(sdkInstance.GetContract())
  if err != nil {
    return nil, fmt.Errorf("AdaptInstance:  %w", err)
  }

  sdkAutoScalingGroup, _ := sdkInstance.GetAutoScalingGroupOk()
  if sdkAutoScalingGroup != nil {
    autoScalingGroup, err = adaptAutoScalingGroup(*sdkAutoScalingGroup)
    if err != nil {
      return nil, fmt.Errorf("AdaptInstance:  %w", err)
    }
  }

  optionalValues := public_cloud.OptionalInstanceValues{
    Reference:        shared.AdaptNullableStringToValue(sdkInstance.Reference),
    MarketAppId:      shared.AdaptNullableStringToValue(sdkInstance.MarketAppId),
    StartedAt:        shared.AdaptNullableTimeToValue(sdkInstance.StartedAt),
    AutoScalingGroup: autoScalingGroup,
  }

  instance := public_cloud.NewInstance(
    *instanceId,
    sdkInstance.GetRegion(),
    adaptResources(sdkInstance.GetResources()),
    adaptImage(sdkInstance.GetImage()),
    state,
    sdkInstance.GetProductType(),
    sdkInstance.GetHasPublicIpV4(),
    sdkInstance.GetIncludesPrivateNetwork(),
    *rootDiskSize,
    public_cloud.InstanceType{Name: string(sdkInstance.GetType())},
    rootDiskStorageType,
    ips,
    *contract,
    optionalValues,
  )

  return &instance, nil
}

func AdaptInstanceDetails(
  sdkInstanceDetails publicCloud.InstanceDetails,
) (*public_cloud.Instance, error) {
  var autoScalingGroup *public_cloud.AutoScalingGroup

  instanceId, err := value_object.NewUuid(sdkInstanceDetails.GetId())
  if err != nil {
    return nil, fmt.Errorf("AdaptInstanceDetails: %w", err)
  }

  state, err := enum.NewState(string(sdkInstanceDetails.GetState()))
  if err != nil {
    return nil, fmt.Errorf("AdaptInstanceDetails: %w", err)
  }

  rootDiskSize, err := value_object.NewRootDiskSize(
    int(sdkInstanceDetails.GetRootDiskSize()),
  )
  if err != nil {
    return nil, fmt.Errorf("AdaptInstanceDetails: %w", err)
  }

  rootDiskStorageType, err := enum.NewRootDiskStorageType(
    string(sdkInstanceDetails.GetRootDiskStorageType()),
  )
  if err != nil {
    return nil, fmt.Errorf("AdaptInstanceDetails: %w", err)
  }

  ips, err := adaptIpsDetails(sdkInstanceDetails.GetIps())
  if err != nil {
    return nil, fmt.Errorf("AdaptInstanceDetails:  %w", err)
  }

  contract, err := adaptContract(sdkInstanceDetails.GetContract())
  if err != nil {
    return nil, fmt.Errorf("AdaptInstanceDetails:  %w", err)
  }

  sdkAutoScalingGroup, _ := sdkInstanceDetails.GetAutoScalingGroupOk()
  if sdkAutoScalingGroup != nil {
    autoScalingGroup, err = adaptAutoScalingGroup(*sdkAutoScalingGroup)
    if err != nil {
      return nil, fmt.Errorf("AdaptInstanceDetails:  %w", err)
    }
  }

  optionalValues := public_cloud.OptionalInstanceValues{
    Reference: shared.AdaptNullableStringToValue(
      sdkInstanceDetails.Reference,
    ),
    MarketAppId: shared.AdaptNullableStringToValue(
      sdkInstanceDetails.MarketAppId,
    ),
    StartedAt:        shared.AdaptNullableTimeToValue(sdkInstanceDetails.StartedAt),
    AutoScalingGroup: autoScalingGroup,
  }
  if sdkInstanceDetails.Iso.Get() != nil {
    iso := adaptIso(*sdkInstanceDetails.Iso.Get())
    optionalValues.Iso = &iso
  }
  if sdkInstanceDetails.PrivateNetwork.Get() != nil {
    privateNetwork := adaptPrivateNetwork(
      *sdkInstanceDetails.PrivateNetwork.Get(),
    )
    optionalValues.PrivateNetwork = &privateNetwork
  }
  if sdkInstanceDetails.Volume.Get() != nil {
    volume := adaptVolume(
      *sdkInstanceDetails.Volume.Get(),
    )
    optionalValues.Volume = &volume
  }

  instance := public_cloud.NewInstance(
    *instanceId,
    sdkInstanceDetails.GetRegion(),
    adaptResources(sdkInstanceDetails.GetResources()),
    adaptImage(sdkInstanceDetails.GetImage()),
    state,
    sdkInstanceDetails.GetProductType(),
    sdkInstanceDetails.GetHasPublicIpV4(),
    sdkInstanceDetails.GetIncludesPrivateNetwork(),
    *rootDiskSize,
    public_cloud.InstanceType{Name: string(sdkInstanceDetails.GetType())},
    rootDiskStorageType,
    ips,
    *contract,
    optionalValues,
  )

  return &instance, nil
}

func adaptResources(sdkResources publicCloud.Resources) public_cloud.Resources {
  resources := public_cloud.NewResources(
    adaptCpu(sdkResources.GetCpu()),
    adaptMemory(sdkResources.GetMemory()),
    adaptNetworkSpeed(sdkResources.GetPublicNetworkSpeed()),
    adaptNetworkSpeed(sdkResources.GetPrivateNetworkSpeed()),
  )

  return resources
}

func adaptCpu(sdkCpu publicCloud.Cpu) public_cloud.Cpu {
  return public_cloud.NewCpu(int(sdkCpu.GetValue()), sdkCpu.GetUnit())
}

func adaptMemory(sdkMemory publicCloud.Memory) public_cloud.Memory {
  return public_cloud.NewMemory(float64(sdkMemory.GetValue()), sdkMemory.GetUnit())
}

func adaptNetworkSpeed(sdkNetworkSpeed publicCloud.NetworkSpeed) public_cloud.NetworkSpeed {
  return public_cloud.NewNetworkSpeed(
    int(sdkNetworkSpeed.GetValue()),
    sdkNetworkSpeed.GetUnit(),
  )
}

func adaptImage(sdkImage publicCloud.Image) public_cloud.Image {
  return public_cloud.NewImage(
    sdkImage.GetId(),
    sdkImage.GetName(),
    nil,
    sdkImage.GetFamily(),
    sdkImage.GetFlavour(),
    nil,
    nil,
    nil,
    nil,
    nil,
    nil,
    sdkImage.GetCustom(),
    nil,
    []string{},
    []string{},
  )
}

func adaptIpsDetails(sdkIps []publicCloud.IpDetails) (public_cloud.Ips, error) {
  var ips public_cloud.Ips
  for _, sdkIp := range sdkIps {
    ip, err := adaptIpDetails(sdkIp)
    if err != nil {
      return nil, fmt.Errorf("adaptIpsDetails: %w", err)
    }
    ips = append(ips, *ip)
  }

  return ips, nil
}

func adaptIps(sdkIps []publicCloud.Ip) (public_cloud.Ips, error) {
  var ips public_cloud.Ips
  for _, sdkIp := range sdkIps {
    ip, err := adaptIp(sdkIp)
    if err != nil {
      return nil, fmt.Errorf("adaptIps: %w", err)
    }
    ips = append(ips, *ip)
  }

  return ips, nil
}

func adaptIpDetails(sdkIp publicCloud.IpDetails) (*public_cloud.Ip, error) {
  networkType, err := enum.NewNetworkType(string(sdkIp.GetNetworkType()))
  if err != nil {
    return nil, fmt.Errorf("adaptIpDetails: %w", err)
  }

  optionalIpValues := public_cloud.OptionalIpValues{
    ReverseLookup: shared.AdaptNullableStringToValue(sdkIp.ReverseLookup),
  }

  sdkDdos, _ := sdkIp.GetDdosOk()
  if sdkDdos != nil {
    ddos := adaptDdos(*sdkDdos)
    optionalIpValues.Ddos = &ddos
  }

  ip := public_cloud.NewIp(
    sdkIp.GetIp(),
    sdkIp.GetPrefixLength(),
    int(sdkIp.GetVersion()),
    sdkIp.GetNullRouted(),
    sdkIp.GetMainIp(),
    networkType,
    optionalIpValues,
  )

  return &ip, nil
}

func adaptIp(sdkIp publicCloud.Ip) (*public_cloud.Ip, error) {
  networkType, err := enum.NewNetworkType(string(sdkIp.GetNetworkType()))
  if err != nil {
    return nil, fmt.Errorf(
      "adaptIpDetails: %w",
      err,
    )
  }

  optionalIpValues := public_cloud.OptionalIpValues{
    ReverseLookup: shared.AdaptNullableStringToValue(sdkIp.ReverseLookup),
  }

  ip := public_cloud.NewIp(
    sdkIp.GetIp(),
    sdkIp.GetPrefixLength(),
    int(sdkIp.GetVersion()),
    sdkIp.GetNullRouted(),
    sdkIp.GetMainIp(),
    networkType,
    optionalIpValues,
  )

  return &ip, nil
}

func adaptDdos(sdkDdos publicCloud.Ddos) public_cloud.Ddos {
  return public_cloud.NewDdos(
    sdkDdos.GetDetectionProfile(),
    sdkDdos.GetProtectionType(),
  )
}

func adaptContract(sdkContract publicCloud.Contract) (*public_cloud.Contract, error) {
  billingFrequency, err := enum.NewContractBillingFrequency(
    int(sdkContract.GetBillingFrequency()),
  )
  if err != nil {
    return nil, fmt.Errorf("adaptContract: %w", err)
  }

  contractTerm, err := enum.NewContractTerm(int(sdkContract.GetTerm()))
  if err != nil {
    return nil, fmt.Errorf("adaptContract: %w", err)
  }

  contractType, err := enum.NewContractType(string(sdkContract.GetType()))
  if err != nil {
    return nil, fmt.Errorf("adaptContract: %w", err)
  }

  contractState, err := enum.NewContractState(string(sdkContract.GetState()))
  if err != nil {
    return nil, fmt.Errorf("adaptContract: %w", err)
  }

  contract, err := public_cloud.NewContract(
    billingFrequency,
    contractTerm,
    contractType,
    sdkContract.GetRenewalsAt(),
    sdkContract.GetCreatedAt(),
    contractState,
    shared.AdaptNullableTimeToValue(sdkContract.EndsAt),
  )

  if err != nil {
    return nil, fmt.Errorf("adaptContract: %w", err)
  }

  return contract, nil
}

func adaptIso(sdkIso publicCloud.Iso) public_cloud.Iso {
  return public_cloud.NewIso(sdkIso.GetId(), sdkIso.GetName())
}

func adaptPrivateNetwork(sdkPrivateNetwork publicCloud.PrivateNetwork) public_cloud.PrivateNetwork {
  return public_cloud.PrivateNetwork{
    Id:     sdkPrivateNetwork.GetPrivateNetworkId(),
    Status: sdkPrivateNetwork.GetStatus(),
    Subnet: sdkPrivateNetwork.GetSubnet(),
  }
}

func AdaptAutoScalingGroupDetails(
  sdkAutoScalingGroup publicCloud.AutoScalingGroupDetails,
) (
  *public_cloud.AutoScalingGroup,
  error,
) {
  var loadBalancer *public_cloud.LoadBalancer

  autoScalingGroupId, err := value_object.NewUuid(sdkAutoScalingGroup.GetId())
  if err != nil {
    return nil, fmt.Errorf("adaptAutoScalingGroupDetails: %w", err)
  }

  autoScalingGroupType, err := enum.NewAutoScalingGroupType(
    string(sdkAutoScalingGroup.GetType()),
  )
  if err != nil {
    return nil, fmt.Errorf("adaptAutoScalingGroupDetails: %w", err)
  }

  state, err := enum.NewAutoScalingGroupState(
    string(sdkAutoScalingGroup.GetState()),
  )
  if err != nil {
    return nil, fmt.Errorf("adaptAutoScalingGroupDetails: %w", err)
  }

  reference, err := value_object.NewAutoScalingGroupReference(
    sdkAutoScalingGroup.GetReference(),
  )
  if err != nil {
    return nil, fmt.Errorf("adaptAutoScalingGroupDetails: %w", err)
  }

  sdkLoadBalancer, _ := sdkAutoScalingGroup.GetLoadBalancerOk()
  if sdkLoadBalancer != nil {
    loadBalancer, err = adaptLoadBalancer(*sdkLoadBalancer)
    if err != nil {
      return nil, fmt.Errorf("adaptAutoScalingGroupDetails: %w", err)
    }
  }

  options := public_cloud.AutoScalingGroupOptions{
    DesiredAmount: shared.AdaptNullableInt32ToValue(sdkAutoScalingGroup.DesiredAmount),
    MinimumAmount: shared.AdaptNullableInt32ToValue(sdkAutoScalingGroup.MinimumAmount),
    MaximumAmount: shared.AdaptNullableInt32ToValue(sdkAutoScalingGroup.MaximumAmount),
    CpuThreshold:  shared.AdaptNullableInt32ToValue(sdkAutoScalingGroup.CpuThreshold),
    CoolDownTime:  shared.AdaptNullableInt32ToValue(sdkAutoScalingGroup.CooldownTime),
    StartsAt:      shared.AdaptNullableTimeToValue(sdkAutoScalingGroup.StartsAt),
    EndsAt:        shared.AdaptNullableTimeToValue(sdkAutoScalingGroup.EndsAt),
    WarmupTime:    shared.AdaptNullableInt32ToValue(sdkAutoScalingGroup.WarmupTime),
    LoadBalancer:  loadBalancer,
  }

  autoScalingGroup := public_cloud.NewAutoScalingGroup(
    *autoScalingGroupId,
    autoScalingGroupType,
    state,
    sdkAutoScalingGroup.GetRegion(),
    *reference,
    sdkAutoScalingGroup.GetCreatedAt(),
    sdkAutoScalingGroup.GetUpdatedAt(),
    options,
  )

  return &autoScalingGroup, nil
}

func AdaptLoadBalancerDetails(
  sdkLoadBalancer publicCloud.LoadBalancerDetails,
) (
  *public_cloud.LoadBalancer,
  error,
) {
  loadBalancerId, err := value_object.NewUuid(sdkLoadBalancer.Id)
  if err != nil {
    return nil, fmt.Errorf("AdaptLoadBalancerDetails: %w", err)
  }

  state, err := enum.NewState(string(sdkLoadBalancer.GetState()))
  if err != nil {
    return nil, fmt.Errorf("AdaptLoadBalancerDetails: %w", err)
  }

  contract, err := adaptContract(sdkLoadBalancer.GetContract())
  if err != nil {
    return nil, fmt.Errorf("AdaptLoadBalancerDetails: %w", err)
  }

  ips, err := adaptIpsDetails(sdkLoadBalancer.GetIps())
  if err != nil {
    return nil, fmt.Errorf("AdaptLoadBalancerDetails:  %w", err)
  }

  options := public_cloud.OptionalLoadBalancerValues{
    Reference: shared.AdaptNullableStringToValue(sdkLoadBalancer.Reference),
    StartedAt: shared.AdaptNullableTimeToValue(sdkLoadBalancer.StartedAt),
  }

  if sdkLoadBalancer.Configuration.Get() != nil {
    configuration, err := adaptLoadBalancerConfiguration(sdkLoadBalancer.GetConfiguration())
    if err != nil {
      return nil, fmt.Errorf("AdaptLoadBalancerDetails:  %w", err)
    }
    options.Configuration = configuration
  }

  if sdkLoadBalancer.PrivateNetwork.Get() != nil {
    privateNetwork := adaptPrivateNetwork(*sdkLoadBalancer.PrivateNetwork.Get())
    options.PrivateNetwork = &privateNetwork
  }

  loadBalancer := public_cloud.NewLoadBalancer(
    *loadBalancerId,
    public_cloud.InstanceType{Name: string(sdkLoadBalancer.GetType())},
    adaptResources(sdkLoadBalancer.GetResources()),
    sdkLoadBalancer.GetRegion(),
    state,
    *contract,
    ips,
    options,
  )

  return &loadBalancer, nil
}

func adaptLoadBalancer(sdkLoadBalancer publicCloud.LoadBalancer) (
  *public_cloud.LoadBalancer,
  error,
) {
  loadBalancerId, err := value_object.NewUuid(sdkLoadBalancer.Id)
  if err != nil {
    return nil, fmt.Errorf("adaptLoadBalancer: %w", err)
  }

  state, err := enum.NewState(string(sdkLoadBalancer.GetState()))
  if err != nil {
    return nil, fmt.Errorf("adaptLoadBalancer: %w", err)
  }

  options := public_cloud.OptionalLoadBalancerValues{
    Reference: shared.AdaptNullableStringToValue(sdkLoadBalancer.Reference),
    StartedAt: shared.AdaptNullableTimeToValue(sdkLoadBalancer.StartedAt),
  }

  loadBalancer := public_cloud.NewLoadBalancer(
    *loadBalancerId,
    public_cloud.InstanceType{Name: string(sdkLoadBalancer.GetType())},
    adaptResources(sdkLoadBalancer.GetResources()),
    "",
    state,
    public_cloud.Contract{},
    public_cloud.Ips{},
    options,
  )

  return &loadBalancer, nil
}

func adaptLoadBalancerConfiguration(sdkLoadBalancerConfiguration publicCloud.LoadBalancerConfiguration) (
  *public_cloud.LoadBalancerConfiguration,
  error,
) {
  balance, err := enum.NewBalance(string(sdkLoadBalancerConfiguration.GetBalance()))
  if err != nil {
    return nil, fmt.Errorf("adaptLoadBalancerConfiguration: %w", err)
  }

  options := public_cloud.OptionalLoadBalancerConfigurationOptions{
    HealthCheck: nil,
  }
  if sdkLoadBalancerConfiguration.StickySession.Get() != nil {
    stickySession := adaptStickySession(*sdkLoadBalancerConfiguration.StickySession.Get())
    options.StickySession = &stickySession
  }
  if sdkLoadBalancerConfiguration.HealthCheck.Get() != nil {
    healthCheck, err := adaptHealthCheck(*sdkLoadBalancerConfiguration.HealthCheck.Get())
    if err != nil {
      return nil, fmt.Errorf("adaptLoadBalancerConfiguration: %w", err)
    }

    options.HealthCheck = healthCheck
  }

  configuration := public_cloud.NewLoadBalancerConfiguration(
    balance,
    sdkLoadBalancerConfiguration.GetXForwardedFor(),
    int(sdkLoadBalancerConfiguration.GetIdleTimeOut()),
    int(sdkLoadBalancerConfiguration.GetTargetPort()),
    options,
  )

  return &configuration, nil
}

func adaptStickySession(sdkStickySession publicCloud.StickySession) public_cloud.StickySession {
  return public_cloud.NewStickySession(
    sdkStickySession.GetEnabled(),
    int(sdkStickySession.GetMaxLifeTime()),
  )
}

func adaptHealthCheck(sdkHealthCheck publicCloud.HealthCheck) (
  *public_cloud.HealthCheck,
  error,
) {
  method, err := enum.NewMethod(sdkHealthCheck.GetMethod())
  if err != nil {
    return nil, fmt.Errorf("adaptHealthCheck: %w", err)
  }

  healthCheck := public_cloud.NewHealthCheck(
    method,
    sdkHealthCheck.GetUri(),
    int(sdkHealthCheck.GetPort()),
    public_cloud.OptionalHealthCheckValues{
      Host: shared.AdaptNullableStringToValue(sdkHealthCheck.Host),
    },
  )

  return &healthCheck, nil
}

func AdaptInstanceType(sdkInstanceType publicCloud.InstanceType) (
  *public_cloud.InstanceType,
  error,
) {
  resources := adaptResources(sdkInstanceType.GetResources())
  prices := adaptPrices(sdkInstanceType.GetPrices())

  optional := public_cloud.OptionalInstanceTypeValues{}

  sdkStorageTypes, _ := sdkInstanceType.GetStorageTypesOk()
  if sdkStorageTypes != nil {
    storageTypes, err := adaptStorageTypes(sdkStorageTypes)
    if err != nil {
      return nil, fmt.Errorf("AdaptInstanceType: %w", err)
    }
    optional.StorageTypes = storageTypes
  }

  instanceType := public_cloud.NewInstanceType(
    sdkInstanceType.GetName(),
    resources,
    prices,
    optional,
  )

  return &instanceType, nil
}

func adaptPrices(sdkPrices publicCloud.Prices) public_cloud.Prices {
  return public_cloud.NewPrices(
    sdkPrices.GetCurrency(),
    sdkPrices.GetCurrencySymbol(),
    adaptPrice(sdkPrices.GetCompute()),
    adaptStorage(sdkPrices.GetStorage()),
  )
}

func adaptStorage(sdkStorage publicCloud.Storage) public_cloud.Storage {
  return public_cloud.NewStorage(
    adaptPrice(sdkStorage.Local),
    adaptPrice(sdkStorage.Central),
  )
}

func adaptPrice(sdkPrice publicCloud.Price) public_cloud.Price {
  return public_cloud.NewPrice(sdkPrice.GetHourlyPrice(), sdkPrice.GetMonthlyPrice())
}

func adaptStorageTypes(sdkStorageTypes []publicCloud.RootDiskStorageType) (
  *public_cloud.StorageTypes,
  error,
) {
  var storageTypes public_cloud.StorageTypes

  for _, sdkStorageType := range sdkStorageTypes {
    storageType, err := enum.NewRootDiskStorageType(string(sdkStorageType))
    if err != nil {
      return nil, fmt.Errorf("adaptStorageTypes: %w", err)
    }
    storageTypes = append(storageTypes, storageType)
  }

  return &storageTypes, nil
}

func AdaptRegion(sdkRegion publicCloud.Region) public_cloud.Region {
  return public_cloud.NewRegion(sdkRegion.GetName(), sdkRegion.GetLocation())
}

func adaptAutoScalingGroup(
  sdkAutoScalingGroup publicCloud.AutoScalingGroup,
) (
  *public_cloud.AutoScalingGroup,
  error,
) {
  autoScalingGroupId, err := value_object.NewUuid(sdkAutoScalingGroup.GetId())
  if err != nil {
    return nil, fmt.Errorf("adaptAutoScalingGroup: %w", err)
  }

  autoScalingGroupType, err := enum.NewAutoScalingGroupType(
    string(sdkAutoScalingGroup.GetType()),
  )
  if err != nil {
    return nil, fmt.Errorf("adaptAutoScalingGroup: %w", err)
  }

  state, err := enum.NewAutoScalingGroupState(string(sdkAutoScalingGroup.GetState()))
  if err != nil {
    return nil, fmt.Errorf("adaptAutoScalingGroup: %w", err)
  }

  reference, err := value_object.NewAutoScalingGroupReference(sdkAutoScalingGroup.GetReference())

  if err != nil {
    return nil, fmt.Errorf("adaptAutoScalingGroup: %w", err)
  }

  options := public_cloud.AutoScalingGroupOptions{
    DesiredAmount: shared.AdaptNullableInt32ToValue(sdkAutoScalingGroup.DesiredAmount),
    MinimumAmount: shared.AdaptNullableInt32ToValue(sdkAutoScalingGroup.MinimumAmount),
    MaximumAmount: shared.AdaptNullableInt32ToValue(sdkAutoScalingGroup.MaximumAmount),
    CpuThreshold:  shared.AdaptNullableInt32ToValue(sdkAutoScalingGroup.CpuThreshold),
    CoolDownTime:  shared.AdaptNullableInt32ToValue(sdkAutoScalingGroup.CooldownTime),
    StartsAt:      shared.AdaptNullableTimeToValue(sdkAutoScalingGroup.StartsAt),
    EndsAt:        shared.AdaptNullableTimeToValue(sdkAutoScalingGroup.EndsAt),
    WarmupTime:    shared.AdaptNullableInt32ToValue(sdkAutoScalingGroup.WarmupTime),
  }

  autoScalingGroup := public_cloud.NewAutoScalingGroup(
    *autoScalingGroupId,
    autoScalingGroupType,
    state,
    sdkAutoScalingGroup.GetRegion(),
    *reference,
    sdkAutoScalingGroup.GetCreatedAt(),
    sdkAutoScalingGroup.GetUpdatedAt(),
    options,
  )

  return &autoScalingGroup, nil
}

func adaptVolume(sdkVolume publicCloud.Volume) public_cloud.Volume {
  return public_cloud.NewVolume(float64(sdkVolume.GetSize()), sdkVolume.GetUnit())
}

func AdaptImageDetails(sdkImageDetails publicCloud.ImageDetails) public_cloud.Image {
  state, _ := sdkImageDetails.GetStateOk()
  stateReason, _ := sdkImageDetails.GetStateReasonOk()
  region, _ := sdkImageDetails.GetRegionOk()
  createdAt, _ := sdkImageDetails.GetCreatedAtOk()
  updatedAt, _ := sdkImageDetails.GetUpdatedAtOk()
  storageSize := adaptStorageSize(sdkImageDetails.GetStorageSize())
  version, _ := sdkImageDetails.GetVersionOk()
  architecture, _ := sdkImageDetails.GetArchitectureOk()

  return public_cloud.NewImage(
    sdkImageDetails.GetId(),
    sdkImageDetails.GetName(),
    version,
    sdkImageDetails.GetFamily(),
    sdkImageDetails.GetFlavour(),
    architecture,
    state,
    stateReason,
    region,
    createdAt,
    updatedAt,
    sdkImageDetails.GetCustom(),
    &storageSize,
    sdkImageDetails.GetMarketApps(),
    sdkImageDetails.GetStorageTypes(),
  )
}

func adaptStorageSize(sdkStorageSize publicCloud.StorageSize) public_cloud.StorageSize {
  return public_cloud.NewStorageSize(
    float64(sdkStorageSize.GetSize()),
    sdkStorageSize.GetUnit(),
  )
}
