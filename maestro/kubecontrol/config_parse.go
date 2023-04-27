package kubecontrol

// TODO Might need this in the future, keeping it for now
//import (
//	"fmt"
//	k8scorev1 "k8s.io/api/core/v1"
//	"k8s.io/apimachinery/pkg/api/resource"
//	k8sv1acapps "k8s.io/client-go/applyconfigurations/apps/v1"
//	k8sv1accore "k8s.io/client-go/applyconfigurations/core/v1"
//	k8sv1acmeta "k8s.io/client-go/applyconfigurations/meta/v1"
//)
//
//func CreateDeploymentApplyConfig(otelConfigYaml, sinkID string) (k8sacv1apps.DeploymentApplyConfiguration, error) {
//	_ := buildConfigMapEntry(otelConfigYaml, sinkID)
//	deployment := buildDeployment(sinkID)
//
//	return deployment, nil
//}
//
//func buildDeployment(sinkID string) k8sacv1apps.DeploymentApplyConfiguration {
//	deploymentMetaName := fmt.Sprintf("otel-%s", sinkID)
//	configMapNameMeta := fmt.Sprintf("otel-collector-config-%s", sinkID)
//	deploymentMetaLabelComponent := fmt.Sprintf("otel-collector-%s", sinkID)
//	labels := map[string]string{
//		"app":       "opentelemetry",
//		"component": deploymentMetaLabelComponent,
//	}
//	deployment := k8sacv1apps.DeploymentApplyConfiguration{
//		ObjectMetaApplyConfiguration: &k8sacv1meta.ObjectMetaApplyConfiguration{
//			Name:              &deploymentMetaName,
//			CreationTimestamp: nil,
//			Labels:            labels,
//		},
//	}
//	deploymentSpec := k8sacv1apps.DeploymentSpec()
//	deploymentSpec.WithReplicas(1)
//	selector := k8sacv1meta.LabelSelector()
//	selector.WithMatchLabels(labels)
//	deploymentSpec.WithSelector(selector)
//	templateSpec := k8sacv1core.PodTemplateSpecApplyConfiguration{
//		ObjectMetaApplyConfiguration: getMetadata(deploymentMetaLabelComponent),
//	}
//	templateSpec.WithLabels(labels)
//	podSpec := k8sacv1core.PodSpec()
//	logVolume := k8sacv1core.Volume()
//	logVolume.WithName("varlog")
//	logVolume.WithHostPath(k8sacv1core.HostPathVolumeSource().WithPath("/var/log").WithType(""))
//	containersVolume := k8sacv1core.Volume()
//	containersVolume.WithName("varlibdockercontainers")
//	containersVolume.WithHostPath(k8sacv1core.HostPathVolumeSource().WithPath("/var/lib/docker/containers").WithType(""))
//	configVolume := k8sacv1core.Volume()
//	configVolume.WithName("data").WithConfigMap(k8sacv1core.ConfigMapVolumeSource().WithName(configMapNameMeta).WithDefaultMode(420))
//	podSpec.WithVolumes(logVolume, containersVolume, configVolume)
//	templateSpec.WithSpec(podSpec)
//	containerPodSpec := k8sacv1core.Container()
//	containerPodSpec.WithName("otel-collector").WithImage("otel/opentelemetry-collector-contrib:0.68.0")
//	heathCheckPort := k8sacv1core.ContainerPort().WithContainerPort(13133).WithProtocol(k8scorev1.ProtocolTCP)
//	pprofPort := k8sacv1core.ContainerPort().WithContainerPort(8888).WithProtocol(k8scorev1.ProtocolTCP)
//	containerPodSpec.WithPorts(heathCheckPort, pprofPort)
//
//	cpuQuantity := resource.NewQuantity(100, "m")
//	memQuantity := resource.NewQuantity(200, "Mi")
//	resourceReqs := k8sacv1core.ResourceRequirements().WithLimits(map[k8scorev1.ResourceName]resource.Quantity{
//		k8scorev1.ResourceCPU:    *cpuQuantity,
//		k8scorev1.ResourceMemory: *memQuantity,
//	}).WithRequests(map[k8scorev1.ResourceName]resource.Quantity{
//		k8scorev1.ResourceCPU:    *cpuQuantity,
//		k8scorev1.ResourceMemory: *memQuantity,
//	})
//	containerPodSpec.WithResources(resourceReqs)
//	logVolumeMount := k8sacv1core.VolumeMount()
//	logVolumeMount.WithName("varlog")
//	containerPodSpec.WithVolumeMounts(k8sacv1core.VolumeMount().)
//	podSpec.WithContainers()
//	deploymentSpec.WithTemplate(&templateSpec)
//	deployment.WithAPIVersion("k8sacv1core")
//	deployment.WithKind("Deployment")
//	deployment.WithSpec()
//	return deployment
//}
//
//func buildConfigMapEntry(otelConfigYaml string, sinkID string) *k8sacv1core.ConfigMapApplyConfiguration {
//	configMapNameMeta := fmt.Sprintf("otel-collector-config-%s", sinkID)
//	configMapDataEntries := make(map[string]string)
//	configMapDataEntries["config.yaml"] = otelConfigYaml
//	metaApplyConfiguration := getMetadata(configMapNameMeta)
//	configMap := k8sacv1core.ConfigMapApplyConfiguration{
//		ObjectMetaApplyConfiguration: metaApplyConfiguration,
//	}
//	configMap.WithKind("ConfigMap")
//	configMap.WithAPIVersion("k8sacv1core")
//	configMap.WithData(configMapDataEntries)
//	return &configMap
//}
//
//func getMetadata(metadaName string) *k8sacv1meta.ObjectMetaApplyConfiguration {
//	metaApplyConfiguration := &k8sacv1meta.ObjectMetaApplyConfiguration{
//		Name:              &metadaName,
//		CreationTimestamp: nil,
//	}
//	return metaApplyConfiguration
//}
