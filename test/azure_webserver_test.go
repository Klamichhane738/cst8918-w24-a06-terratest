package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

var subscriptionID string = "b296b604-9c48-4e98-bb66-56c661c39a1d"

func TestAzureLinuxVMCreation(t *testing.T) {
	start := time.Now()

	terraformOptions := &terraform.Options{
		TerraformDir: "../",
		Vars: map[string]interface{}{
			"labelPrefix": "lami0053",
		},
	}

	terraform.InitAndApply(t, terraformOptions)

	vmName := terraform.Output(t, terraformOptions, "vm_name")
	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")
	nicName := terraform.Output(t, terraformOptions, "nic_name") // Assuming you output the NIC name

	// Expected Ubuntu image details
	expectedPublisher := "Canonical"
	expectedOffer := "0001-com-ubuntu-server-jammy"
	expectedSku := "22_04-lts-gen2"
	expectedVersion := "latest" // Or you could check for the specific version if needed

	// Check if the virtual machine exists
	assert.True(t, azure.VirtualMachineExists(t, vmName, resourceGroupName, subscriptionID))

	// Confirm NIC exists and is attached to the VM
	nicList := azure.GetVirtualMachineNics(t, vmName, resourceGroupName, subscriptionID)
	nicExists := false
	for _, nic := range nicList {
		if nic == nicName {
			nicExists = true
			break
		}
	}
	assert.True(t, nicExists, "Network interface %s is not attached to VM %s", nicName, vmName)

	// Confirm that the VM is running the correct Ubuntu version
	// Get the VM image details to check the OS version
	vmImage := azure.GetVirtualMachineImage(t, vmName, resourceGroupName, subscriptionID)
	assert.Equal(t, expectedPublisher, vmImage.Publisher, "VM Publisher does not match expected value")
	assert.Equal(t, expectedOffer, vmImage.Offer, "VM Offer does not match expected value")
	assert.Equal(t, expectedSku, vmImage.SKU, "VM SKU does not match expected value")
	assert.Equal(t, expectedVersion, vmImage.Version, "Ubuntu version does not match expected version")

	t.Log("PASS")
	duration := time.Since(start)
	fmt.Printf("--- PASS: TestAzureLinuxVMCreation (%.2fs)\n", duration.Seconds())
}
