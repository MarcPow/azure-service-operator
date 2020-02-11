// +build all eventhub onlyeventhub

package controllers

import (
	"context"
	"strings"
	"testing"

	azurev1alpha1 "github.com/Azure/azure-service-operator/api/v1alpha1"
	"github.com/Azure/azure-service-operator/pkg/helpers"
	"github.com/Azure/azure-service-operator/pkg/resourcemanager/config"
	kvsecrets "github.com/Azure/azure-service-operator/pkg/secrets/keyvault"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
)

/*
func TestEventHubControllerNoNamespace(t *testing.T) {
	t.Parallel()
	defer PanicRecover()
	ctx := context.Background()
	assert := assert.New(t)

	// Add Tests for OpenAPI validation (or additonal CRD features) specified in
	// your API definition.
	// Avoid adding tests for vanilla CRUD operations because they would
	// test Kubernetes API server, which isn't the goal here.

	eventhubName := "t-eh-" + helpers.RandomString(10)

	// Create the EventHub object and expect the Reconcile to be created
	eventhubInstance := &azurev1alpha1.Eventhub{
		ObjectMeta: metav1.ObjectMeta{
			Name:      eventhubName,
			Namespace: "default",
		},
		Spec: azurev1alpha1.EventhubSpec{
			Location:      "westus",
			Namespace:     "t-ns-dev-eh-" + helpers.RandomString(10),
			ResourceGroup: tc.resourceGroupName,
			Properties: azurev1alpha1.EventhubProperties{
				MessageRetentionInDays: 7,
				PartitionCount:         2,
			},
		},
	}

	err := tc.k8sClient.Create(ctx, eventhubInstance)
	assert.Equal(nil, err, "create eventhub in k8s")

	eventhubNamespacedName := types.NamespacedName{Name: eventhubName, Namespace: "default"}

	assert.Eventually(func() bool {
		_ = tc.k8sClient.Get(ctx, eventhubNamespacedName, eventhubInstance)
		return strings.Contains(eventhubInstance.Status.Message, errhelp.ParentNotFoundErrorCode)
	}, tc.timeout, tc.retry, "wait for eventhub to provision")

	err = tc.k8sClient.Delete(ctx, eventhubInstance)
	assert.Equal(nil, err, "delete eventhub in k8s")

	assert.Eventually(func() bool {
		err = tc.k8sClient.Get(ctx, eventhubNamespacedName, eventhubInstance)
		return apierrors.IsNotFound(err)
	}, tc.timeout, tc.retry, "wait for eventHubInstance to be gone from k8s")

}

func TestEventHubControllerCreateAndDelete(t *testing.T) {
	t.Parallel()
	defer PanicRecover()
	ctx := context.Background()
	assert := assert.New(t)

	// Add any setup steps that needs to be executed before each test
	rgName := tc.resourceGroupName
	ehnName := tc.eventhubNamespaceName
	eventhubName := "t-eh-" + helpers.RandomString(10)

	// Create the EventHub object and expect the Reconcile to be created
	eventhubInstance := &azurev1alpha1.Eventhub{
		ObjectMeta: metav1.ObjectMeta{
			Name:      eventhubName,
			Namespace: "default",
		},
		Spec: azurev1alpha1.EventhubSpec{
			Location:      "westus",
			Namespace:     ehnName,
			ResourceGroup: rgName,
			Properties: azurev1alpha1.EventhubProperties{
				MessageRetentionInDays: 7,
				PartitionCount:         2,
			},
			AuthorizationRule: azurev1alpha1.EventhubAuthorizationRule{
				Name:   "RootManageSharedAccessKey",
				Rights: []string{"Listen"},
			},
		},
	}

	err := tc.k8sClient.Create(ctx, eventhubInstance)
	assert.Equal(nil, err, "create eventhub in k8s")

	eventhubNamespacedName := types.NamespacedName{Name: eventhubName, Namespace: "default"}

	assert.Eventually(func() bool {
		_ = tc.k8sClient.Get(ctx, eventhubNamespacedName, eventhubInstance)
		return eventhubInstance.HasFinalizer(finalizerName)
	}, tc.timeout, tc.retry, "wait for eventhub to have finalizer")

	assert.Eventually(func() bool {
		_ = tc.k8sClient.Get(ctx, eventhubNamespacedName, eventhubInstance)
		return strings.Contains(eventhubInstance.Status.Message, successMsg)
	}, tc.timeout, tc.retry, "wait for eventhub to provision")

	err = tc.k8sClient.Delete(ctx, eventhubInstance)
	assert.Equal(nil, err, "delete eventhub in k8s")

	assert.Eventually(func() bool {
		err = tc.k8sClient.Get(ctx, eventhubNamespacedName, eventhubInstance)
		return apierrors.IsNotFound(err)
	}, tc.timeout, tc.retry, "wait for eventHubInstance to be gone from k8s")

}

func TestEventHubControllerCreateAndDeleteCustomSecret(t *testing.T) {
	t.Parallel()
	defer PanicRecover()
	ctx := context.Background()
	assert := assert.New(t)

	// Add any setup steps that needs to be executed before each test
	rgName := tc.resourceGroupName
	rgLocation := tc.resourceGroupLocation
	ehnName := tc.eventhubNamespaceName
	eventhubName := "t-eh-" + helpers.RandomString(10)
	secretName := "secret-" + eventhubName

	// Create the EventHub object and expect the Reconcile to be created
	eventhubInstance := &azurev1alpha1.Eventhub{
		ObjectMeta: metav1.ObjectMeta{
			Name:      eventhubName,
			Namespace: "default",
		},
		Spec: azurev1alpha1.EventhubSpec{
			Location:      rgLocation,
			Namespace:     ehnName,
			ResourceGroup: rgName,
			Properties: azurev1alpha1.EventhubProperties{
				MessageRetentionInDays: 7,
				PartitionCount:         2,
			},
			AuthorizationRule: azurev1alpha1.EventhubAuthorizationRule{
				Name:   "RootManageSharedAccessKey",
				Rights: []string{"Listen"},
			},
			SecretName: secretName,
		},
	}

	err := tc.k8sClient.Create(ctx, eventhubInstance)
	assert.Equal(nil, err, "create eventhub in k8s")

	eventhubNamespacedName := types.NamespacedName{Name: eventhubName, Namespace: "default"}

	assert.Eventually(func() bool {
		_ = tc.k8sClient.Get(ctx, eventhubNamespacedName, eventhubInstance)
		return eventhubInstance.HasFinalizer(finalizerName)
	}, tc.timeout, tc.retry, "wait for eventhub to have finalizer")

	assert.Eventually(func() bool {
		_ = tc.k8sClient.Get(ctx, eventhubNamespacedName, eventhubInstance)
		return strings.Contains(eventhubInstance.Status.Message, successMsg)
	}, tc.timeout, tc.retry, "wait for eventhub to provision")

	err = tc.k8sClient.Delete(ctx, eventhubInstance)
	assert.Equal(nil, err, "delete eventhub in k8s")

	assert.Eventually(func() bool {
		err = tc.k8sClient.Get(ctx, eventhubNamespacedName, eventhubInstance)
		return apierrors.IsNotFound(err)
	}, tc.timeout, tc.retry, "wait for eventHubInstance to be gone from k8s")

}
*/
func TestEventHubControllerCreateAndDeleteCustomKeyVault(t *testing.T) {
	t.Parallel()
	defer PanicRecover()
	ctx := context.Background()
	assert := assert.New(t)

	// Add any setup steps that needs to be executed before each test
	rgName := tc.resourceGroupName
	rgLocation := tc.resourceGroupLocation
	ehnName := tc.eventhubNamespaceName
	eventhubName := "t-eh-" + helpers.RandomString(10)
	secretName := "secret-" + eventhubName
	keyVaultNameForSecrets := "t-eh-kv-secrets-" + helpers.RandomString(5)

	allPermissions := []string{"Get", "List", "Set", "Delete", "Recover", "Backup", "Restore"}
	onlyListPermission := []string{"List"}
	// Create the keyVault and add access policies
	keyVaultInstance := &azurev1alpha1.KeyVault{
		ObjectMeta: metav1.ObjectMeta{
			Name:      keyVaultNameForSecrets,
			Namespace: "default",
		},
		Spec: azurev1alpha1.KeyVaultSpec{
			Location:      rgLocation,
			ResourceGroup: rgName,
			AccessPolicies: &[]azurev1alpha1.AccessPolicyEntry{
				{
					TenantID: config.TenantID(),
					ObjectID: config.ClientID(),
					Permissions: &azurev1alpha1.Permissions{
						Keys:         &onlyListPermission,
						Secrets:      &allPermissions,
						Certificates: &onlyListPermission,
						Storage:      &onlyListPermission,
					},
				},
			},
		},
	}

	// Create the Keyvault object and expect the Reconcile to be created
	err := tc.k8sClient.Create(ctx, keyVaultInstance)
	assert.Equal(nil, err, "create keyvault in k8s")

	// Prep query for get
	keyVaultNamespacedName := types.NamespacedName{Name: keyVaultNameForSecrets, Namespace: "default"}

	assert.Eventually(func() bool {
		_ = tc.k8sClient.Get(ctx, keyVaultNamespacedName, keyVaultInstance)
		return helpers.HasFinalizer(keyVaultInstance, finalizerName)
	}, tc.timeout, tc.retry, "wait for keyvault to have finalizer")

	// Wait until key vault is provisioned

	assert.Eventually(func() bool {
		_ = tc.k8sClient.Get(ctx, keyVaultNamespacedName, keyVaultInstance)
		return strings.Contains(keyVaultInstance.Status.Message, successMsg)
	}, tc.timeout, tc.retry, "wait for keyVaultInstance to be ready in k8s")

	// Create the EventHub object and expect the Reconcile to be created
	eventhubInstance := &azurev1alpha1.Eventhub{
		ObjectMeta: metav1.ObjectMeta{
			Name:      eventhubName,
			Namespace: "default",
		},
		Spec: azurev1alpha1.EventhubSpec{
			Location:      rgLocation,
			Namespace:     ehnName,
			ResourceGroup: rgName,
			Properties: azurev1alpha1.EventhubProperties{
				MessageRetentionInDays: 7,
				PartitionCount:         2,
			},
			AuthorizationRule: azurev1alpha1.EventhubAuthorizationRule{
				Name:   "RootManageSharedAccessKey",
				Rights: []string{"Listen"},
			},
			SecretName:             secretName,
			KeyVaultToStoreSecrets: keyVaultNameForSecrets,
		},
	}

	err = tc.k8sClient.Create(ctx, eventhubInstance)
	assert.Equal(nil, err, "create eventhub in k8s")

	eventhubNamespacedName := types.NamespacedName{Name: eventhubName, Namespace: "default"}

	assert.Eventually(func() bool {
		_ = tc.k8sClient.Get(ctx, eventhubNamespacedName, eventhubInstance)
		return eventhubInstance.HasFinalizer(finalizerName)
	}, tc.timeout, tc.retry, "wait for eventhub to have finalizer")

	assert.Eventually(func() bool {
		_ = tc.k8sClient.Get(ctx, eventhubNamespacedName, eventhubInstance)
		return strings.Contains(eventhubInstance.Status.Message, successMsg)
	}, tc.timeout, tc.retry, "wait for eventhub to provision")

	// Check that the secret is added to KeyVault
	keyvaultSecretClient := kvsecrets.New(keyVaultNameForSecrets)
	_, err = keyvaultSecretClient.Get(ctx, keyVaultNamespacedName)
	assert.Equal(nil, err, "checking ßif secret is present in keyvault")

	err = tc.k8sClient.Delete(ctx, eventhubInstance)
	assert.Equal(nil, err, "delete eventhub in k8s")

	assert.Eventually(func() bool {
		err = tc.k8sClient.Get(ctx, eventhubNamespacedName, eventhubInstance)
		return apierrors.IsNotFound(err)
	}, tc.timeout, tc.retry, "wait for eventHubInstance to be gone from k8s")

}

/*
func TestEventHubControllerCreateAndDeleteCapture(t *testing.T) {
	t.Parallel()
	defer PanicRecover()
	ctx := context.Background()
	assert := assert.New(t)

	// Add any setup steps that needs to be executed before each test
	rgName := tc.resourceGroupName
	rgLocation := tc.resourceGroupLocation
	ehnName := tc.eventhubNamespaceName
	saName := tc.storageAccountName
	bcName := tc.blobContainerName
	eventHubName := "t-eh-" + helpers.RandomString(10)

	// Create the EventHub object and expect the Reconcile to be created
	eventhubInstance := &azurev1alpha1.Eventhub{
		ObjectMeta: metav1.ObjectMeta{
			Name:      eventHubName,
			Namespace: "default",
		},
		Spec: azurev1alpha1.EventhubSpec{
			Location:      rgLocation,
			Namespace:     ehnName,
			ResourceGroup: rgName,
			Properties: azurev1alpha1.EventhubProperties{
				MessageRetentionInDays: 7,
				PartitionCount:         2,
				CaptureDescription: azurev1alpha1.CaptureDescription{
					Destination: azurev1alpha1.Destination{
						ArchiveNameFormat: "{Namespace}/{EventHub}/{PartitionId}/{Year}/{Month}/{Day}/{Hour}/{Minute}/{Second}",
						BlobContainer:     bcName,
						Name:              "EventHubArchive.AzureBlockBlob",
						StorageAccount: azurev1alpha1.StorageAccount{
							ResourceGroup: rgName,
							AccountName:   saName,
						},
					},
					Enabled:           true,
					SizeLimitInBytes:  524288000,
					IntervalInSeconds: 300,
				},
			},
			AuthorizationRule: azurev1alpha1.EventhubAuthorizationRule{
				Name: "RootManageSharedAccessKey",
				Rights: []string{
					"Listen",
					"Manage",
					"Send",
				},
			},
		},
	}

	err := tc.k8sClient.Create(ctx, eventhubInstance)
	assert.Equal(nil, err, "create eventhub in k8s")

	eventHubNamespacedName := types.NamespacedName{Name: eventHubName, Namespace: "default"}

	assert.Eventually(func() bool {
		_ = tc.k8sClient.Get(ctx, eventHubNamespacedName, eventhubInstance)
		return eventhubInstance.HasFinalizer(finalizerName)
	}, tc.timeout, tc.retry, "wait for eventhub to have finalizer")

	assert.Eventually(func() bool {
		_ = tc.k8sClient.Get(ctx, eventHubNamespacedName, eventhubInstance)
		return strings.Contains(eventhubInstance.Status.Message, successMsg)
	}, tc.timeout, tc.retry, "wait for eventhub to provision")

	assert.Eventually(func() bool {
		hub, _ := tc.eventhubClient.GetHub(ctx, rgName, ehnName, eventHubName)
		if hub.Properties == nil || hub.CaptureDescription == nil || hub.CaptureDescription.Enabled == nil {
			return false
		}
		return *hub.CaptureDescription.Enabled
	}, tc.timeout, tc.retry, "wait for eventhub capture check")

	err = tc.k8sClient.Delete(ctx, eventhubInstance)
	assert.Equal(nil, err, "delete eventhub in k8s")

	assert.Eventually(func() bool {
		err = tc.k8sClient.Get(ctx, eventHubNamespacedName, eventhubInstance)
		return apierrors.IsNotFound(err)
	}, tc.timeout, tc.retry, "wait for eventHubInstance to be gone from k8s")

}*/
