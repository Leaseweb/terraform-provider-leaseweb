package instances

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
)

func (d *instancesDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state instancesDataSourceModel

	instances, _, err := d.client.SdkClient.PublicCloudAPI.GetInstanceList(d.client.AuthContext()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Leaseweb Instances",
			err.Error(),
		)
		return
	}

	for _, instance := range instances.Instances {
		instanceState := instancesModel{}

		if instance.HasId() {
			instanceState.Id = types.StringValue(instance.GetId())
		}
		if instance.HasEquipmentId() {
			instanceState.EquipmentId = types.StringValue(instance.GetEquipmentId())
		}
		if instance.HasSalesOrgId() {
			instanceState.SalesOrgId = types.StringValue(instance.GetSalesOrgId())
		}
		if instance.HasCustomerId() {
			instanceState.CustomerId = types.StringValue(instance.GetCustomerId())
		}
		if instance.HasRegion() {
			instanceState.Region = types.StringValue(instance.GetRegion())
		}
		if instance.HasReference() {
			instanceState.Reference = types.StringValue(instance.GetReference())
		}
		if instance.HasState() {
			instanceState.State = types.StringValue(string(instance.GetState()))
		}
		if instance.HasProductType() {
			instanceState.ProductType = types.StringValue(instance.GetProductType())
		}
		if instance.HasHasPublicIpV4() {
			instanceState.HasPublicIpv4 = types.BoolValue(instance.GetHasPublicIpV4())
		}
		if instance.HasincludesPrivateNetwork() {
			instanceState.HasPrivateNetwork = types.BoolValue(instance.GetincludesPrivateNetwork())
		}
		if instance.HasType() {
			instanceState.Type = types.StringValue(instance.GetType())
		}
		if instance.HasRootDiskSize() {
			instanceState.RootDiskSize = types.Int64Value(int64(instance.GetRootDiskSize()))
		}
		if instance.HasRootDiskStorageType() {
			instanceState.RootDiskStorageType = types.StringValue(instance.GetRootDiskStorageType())
		}
		if !instance.GetStartedAt().IsZero() {
			instanceState.StartedAt = types.StringValue(instance.GetStartedAt().String())
		}
		if instance.HasMarketAppId() {
			instanceState.MarketAppId = types.StringValue(instance.GetMarketAppId())
		}

		if instance.HasResources() {
			d.setResources(&instanceState, instance.Resources)
		}
		if instance.HasOperatingSystem() {
			d.setOperatingSystem(&instanceState, instance.OperatingSystem)
		}
		if instance.HasIps() {
			d.setIps(&instanceState, instance.Ips)
		}
		if instance.HasContract() {
			d.setContract(&instanceState, instance.Contract)
		}
		if instance.HasIso() {
			d.setIso(&instanceState, instance.Iso.Get())
		}
		if instance.HasPrivateNetwork() {
			d.setPrivateNetwork(&instanceState, instance.PrivateNetwork.Get())
		}

		state.Instances = append(state.Instances, instanceState)
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *instancesDataSource) setPrivateNetwork(instanceState *instancesModel, privateNetwork *publicCloud.PrivateNetwork) {
	privateNetworkState := privateNetworkModel{}

	if privateNetwork.HasPrivateNetworkId() {
		privateNetworkState.Id = types.StringValue(privateNetwork.GetPrivateNetworkId())
	}
	if privateNetwork.HasStatus() {
		privateNetworkState.Status = types.StringValue(privateNetwork.GetStatus())
	}
	if privateNetwork.HasSubnet() {
		privateNetworkState.Subnet = types.StringValue(privateNetwork.GetSubnet())
	}

	instanceState.PrivateNetwork = privateNetworkState
}

func (d *instancesDataSource) setIso(instanceState *instancesModel, iso *publicCloud.Iso) {
	isoState := isoModel{}

	if iso.HasId() {
		isoState.Id = types.StringValue(iso.GetId())
	}
	if iso.HasName() {
		isoState.Name = types.StringValue(iso.GetName())
	}

	instanceState.Iso = isoState
}

func (d *instancesDataSource) setResources(instanceState *instancesModel, resources *publicCloud.InstanceResources) {
	resourcesState := resourcesModel{
		PrivateNetworkSpeed: networkSpeedModel{
			Value: types.Int64Value(int64(resources.PublicNetworkSpeed.GetValue())),
			Unit:  types.StringValue(resources.PublicNetworkSpeed.GetUnit()),
		},
	}

	if resources.HasCpu() {
		cpuState := cpuModel{}

		if resources.Cpu.HasValue() {
			cpuState.Value = types.Int64Value(int64(resources.Cpu.GetValue()))
		}
		if resources.Cpu.HasUnit() {
			cpuState.Unit = types.StringValue(resources.Cpu.GetUnit())
		}

		resourcesState.Cpu = cpuState
	}

	if resources.HasMemory() {
		memoryState := memoryModel{}

		if resources.Memory.HasValue() {
			memoryState.Value = types.Float64Value(float64(resources.Memory.GetValue()))
		}
		if resources.Memory.HasUnit() {
			memoryState.Unit = types.StringValue(resources.Memory.GetUnit())
		}

		resourcesState.Memory = memoryState
	}

	if resources.HasPublicNetworkSpeed() {
		publicNetworkSpeedState := networkSpeedModel{}

		if resources.PublicNetworkSpeed.HasValue() {
			publicNetworkSpeedState.Value = types.Int64Value(int64(resources.PublicNetworkSpeed.GetValue()))
		}
		if resources.PublicNetworkSpeed.HasUnit() {
			publicNetworkSpeedState.Unit = types.StringValue(resources.PublicNetworkSpeed.GetUnit())
		}

		resourcesState.PublicNetworkSpeed = publicNetworkSpeedState
	}

	if resources.HasPrivateNetworkSpeed() {
		privateNetworkSpeedState := networkSpeedModel{}

		if resources.PrivateNetworkSpeed.HasValue() {
			privateNetworkSpeedState.Value = types.Int64Value(int64(resources.PrivateNetworkSpeed.GetValue()))
		}
		if resources.PrivateNetworkSpeed.HasUnit() {
			privateNetworkSpeedState.Unit = types.StringValue(resources.PrivateNetworkSpeed.GetUnit())
		}

		resourcesState.PrivateNetworkSpeed = privateNetworkSpeedState
	}

	instanceState.Resources = resourcesState
}

func (d *instancesDataSource) setOperatingSystem(instanceState *instancesModel, operatingSystem *publicCloud.OperatingSystem) {
	operatingSystemState := operatingSystemModel{}

	if operatingSystem.HasId() {
		operatingSystemState.Id = types.StringValue(operatingSystem.GetId())
	}
	if operatingSystem.HasName() {
		operatingSystemState.Name = types.StringValue(operatingSystem.GetName())
	}
	if operatingSystem.HasFamily() {
		operatingSystemState.Family = types.StringValue(operatingSystem.GetFamily())
	}
	if operatingSystem.HasVersion() {
		operatingSystemState.Version = types.StringValue(operatingSystem.GetVersion())
	}
	if operatingSystem.HasFlavour() {
		operatingSystemState.Flavour = types.StringValue(operatingSystem.GetFlavour())
	}
	if operatingSystem.HasArchitecture() {
		operatingSystemState.Architecture = types.StringValue(operatingSystem.GetArchitecture())
	}

	for _, marketApp := range operatingSystem.MarketApps {
		operatingSystemState.MarketApps = append(
			operatingSystemState.MarketApps, types.StringValue(marketApp),
		)
	}

	for _, storageType := range operatingSystem.StorageTypes {
		operatingSystemState.StorageTypes = append(
			operatingSystemState.StorageTypes, types.StringValue(storageType),
		)
	}

	instanceState.OperatingSystem = operatingSystemState
}

func (d *instancesDataSource) setIps(instanceState *instancesModel, ips []publicCloud.Ip) {
	for _, ip := range ips {
		ipModel := ipModel{}

		if ip.HasIp() {
			ipModel.Ip = types.StringValue(ip.GetIp())
		}
		if ip.HasPrefixLength() {
			ipModel.PrefixLength = types.StringValue(ip.GetPrefixLength())
		}
		if ip.HasVersion() {
			ipModel.Version = types.Int64Value(int64(ip.GetVersion()))
		}
		if ip.HasNullRouted() {
			ipModel.NullRouted = types.BoolValue(ip.GetNullRouted())
		}
		if ip.HasMainIp() {
			ipModel.MainIp = types.BoolValue(ip.GetMainIp())
		}
		if ip.HasNetworkType() {
			ipModel.NetworkType = types.StringValue(string(ip.GetNetworkType()))
		}
		if ip.HasReverseLookup() {
			ipModel.ReverseLookup = types.StringValue(ip.GetReverseLookup())
		}
		if ip.HasDdos() {
			ipModel.Ddos = dDosModel{}

			if ip.Ddos.HasDetectionProfile() {
				ipModel.Ddos.DetectionProfile = types.StringValue(ip.Ddos.GetDetectionProfile())
			}
			if ip.Ddos.HasProtectionType() {
				ipModel.Ddos.ProtectionType = types.StringValue(ip.Ddos.GetProtectionType())
			}
		}

		instanceState.Ips = append(instanceState.Ips, ipModel)
	}
}

func (d *instancesDataSource) setContract(instanceState *instancesModel, contract *publicCloud.Contract) {
	contractState := contractModel{}

	if contract.HasBillingFrequency() {
		contractState.BillingFrequency = types.Int64Value(int64(contract.GetBillingFrequency()))
	}
	if contract.HasTerm() {
		contractState.Term = types.Int64Value(int64(contract.GetTerm()))
	}
	if contract.HasType() {
		contractState.Type = types.StringValue(string(contract.GetType()))
	}
	if !contract.GetEndsAt().IsZero() {
		contractState.EndsAt = types.StringValue(contract.GetEndsAt().String())
	}
	if !contract.GetRenewalsAt().IsZero() {
		contractState.RenewalsAt = types.StringValue(contract.GetRenewalsAt().String())
	}
	if !contract.GetCreatedAt().IsZero() {
		contractState.CreatedAt = types.StringValue(contract.GetCreatedAt().String())
	}
	if contract.HasState() {
		contractState.State = types.StringValue(string(contract.GetState()))
	}

	instanceState.Contract = contractState
}
