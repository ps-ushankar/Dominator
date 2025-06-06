package client

import (
	"net"
	"time"

	"github.com/Cloud-Foundations/Dominator/lib/filesystem"
	"github.com/Cloud-Foundations/Dominator/lib/filter"
	"github.com/Cloud-Foundations/Dominator/lib/log"
	"github.com/Cloud-Foundations/Dominator/lib/srpc"
	proto "github.com/Cloud-Foundations/Dominator/proto/hypervisor"
)

func AcknowledgeVm(client *srpc.Client, ipAddress net.IP) error {
	return acknowledgeVm(client, ipAddress)
}

func AddVmVolumes(client *srpc.Client, ipAddress net.IP, sizes []uint64) error {
	return addVmVolumes(client, ipAddress, sizes)
}

func ChangeVmConsoleType(client *srpc.Client, ipAddress net.IP,
	consoleType proto.ConsoleType) error {
	return changeVmConsoleType(client, ipAddress, consoleType)
}

func ChangeVmCpuPriority(client *srpc.Client, ipAddress net.IP,
	request proto.ChangeVmCpuPriorityRequest) error {
	return changeVmCpuPriority(client, ipAddress, request)
}

func ChangeVmMachineType(client *srpc.Client, ipAddress net.IP,
	machineType proto.MachineType) error {
	return changeVmMachineType(client, ipAddress, machineType)
}

func ChangeVmSize(client *srpc.Client,
	request proto.ChangeVmSizeRequest) error {
	return changeVmSize(client, request)
}

func ChangeVmSubnet(client *srpc.Client,
	request proto.ChangeVmSubnetRequest) (proto.ChangeVmSubnetResponse, error) {
	return changeVmSubnet(client, request)
}

func ChangeVmVolumeInterfaces(client *srpc.Client, ipAddress net.IP,
	volumeInterfaces []proto.VolumeInterface) error {
	return changeVmVolumeInterfaces(client, ipAddress, volumeInterfaces)
}

func ChangeVmVolumeSize(client *srpc.Client, ipAddress net.IP, index uint,
	size uint64) error {
	return changeVmVolumeSize(client, ipAddress, index, size)
}

func ConnectToVmConsole(client *srpc.Client, ipAddr net.IP,
	vncViewerCommand string, logger log.DebugLogger) error {
	return connectToVmConsole(client, ipAddr, vncViewerCommand, logger)
}

func CreateVm(client *srpc.Client, request proto.CreateVmRequest,
	reply *proto.CreateVmResponse, logger log.DebugLogger) error {
	return createVm(client, request, reply, logger)
}

func DeleteVmVolume(client *srpc.Client, ipAddr net.IP, accessToken []byte,
	volumeIndex uint) error {
	return deleteVmVolume(client, ipAddr, accessToken, volumeIndex)
}

func DestroyVm(client *srpc.Client, ipAddr net.IP, accessToken []byte) error {
	return destroyVm(client, ipAddr, accessToken)
}

func ExportLocalVm(client *srpc.Client, ipAddr net.IP,
	verificationCookie []byte) (proto.ExportLocalVmInfo, error) {
	return exportLocalVm(client, ipAddr, verificationCookie)
}

func GetCapacity(client *srpc.Client) (proto.GetCapacityResponse, error) {
	return getCapacity(client)
}

// GetIdentityProvider will get the base URL of the Identity Provider.
func GetIdentityProvider(client srpc.ClientI) (string, error) {
	return getIdentityProvider(client)
}

// GetPublicKey will get the PEM-encoded public key for the Hypervisor.
func GetPublicKey(client srpc.ClientI) ([]byte, error) {
	return getPublicKey(client)
}

func GetRootCookiePath(client *srpc.Client) (string, error) {
	return getRootCookiePath(client)
}

func GetVmInfo(client *srpc.Client, ipAddr net.IP) (proto.VmInfo, error) {
	return getVmInfo(client, ipAddr)
}

func GetVmInfos(client *srpc.Client,
	request proto.GetVmInfosRequest) ([]proto.VmInfo, error) {
	return getVmInfos(client, request)
}

func GetVmLastPatchLog(client *srpc.Client, ipAddr net.IP) (
	[]byte, time.Time, error) {
	return getVmLastPatchLog(client, ipAddr)
}

func HoldLock(client *srpc.Client, timeout time.Duration,
	writeLock bool) error {
	return holdLock(client, timeout, writeLock)
}

func HoldVmLock(client *srpc.Client, ipAddr net.IP, timeout time.Duration,
	writeLock bool) error {
	return holdVmLock(client, ipAddr, timeout, writeLock)
}

func ListSubnets(client *srpc.Client, doSort bool) ([]proto.Subnet, error) {
	return listSubnets(client, doSort)
}

func ListVMs(client *srpc.Client,
	request proto.ListVMsRequest) ([]net.IP, error) {
	return listVMs(client, request)
}

func ListVolumeDirectories(client *srpc.Client, doSort bool) ([]string, error) {
	return listVolumeDirectories(client, doSort)
}

func PowerOff(client *srpc.Client, stopVMs bool) error {
	return powerOff(client, stopVMs)
}

func PrepareVmForMigration(client *srpc.Client, ipAddr net.IP,
	accessToken []byte, enable bool) error {
	return prepareVmForMigration(client, ipAddr, accessToken, enable)
}

func RegisterExternalLeases(client *srpc.Client, addressList proto.AddressList,
	hostnames []string) error {
	return registerExternalLeases(client, addressList, hostnames)
}

func ReorderVmVolumes(client *srpc.Client, ipAddr net.IP, accessToken []byte,
	volumeIndices []uint) error {
	return reorderVmVolumes(client, ipAddr, accessToken, volumeIndices)
}

func ReplaceVmIdentity(client srpc.ClientI,
	request proto.ReplaceVmIdentityRequest) error {
	return replaceVmIdentity(client, request)
}

func ScanVmRoot(client *srpc.Client, ipAddr net.IP,
	scanFilter *filter.Filter) (*filesystem.FileSystem, error) {
	return scanVmRoot(client, ipAddr, scanFilter)
}

func SetDisabledState(client *srpc.Client, disable bool) error {
	return setDisabledState(client, disable)
}

func StartVm(client *srpc.Client, ipAddr net.IP, accessToken []byte) error {
	return startVm(client, ipAddr, accessToken)
}

func StopVm(client *srpc.Client, ipAddr net.IP, accessToken []byte) error {
	return stopVm(client, ipAddr, accessToken)
}
