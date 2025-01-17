package services

import (
	"context"
	"fmt"
	"net"

	apiv1 "github.com/acorn-io/acorn/pkg/apis/api.acorn.io/v1"
	v1 "github.com/acorn-io/acorn/pkg/apis/internal.acorn.io/v1"
	"github.com/acorn-io/acorn/pkg/config"
	"github.com/acorn-io/acorn/pkg/labels"
	"github.com/acorn-io/acorn/pkg/ports"
	"github.com/acorn-io/baaah/pkg/router"
	"github.com/acorn-io/baaah/pkg/typed"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func toContainerLabelsService(service *v1.ServiceInstance) (result []kclient.Object) {
	newService := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        service.Name,
			Namespace:   service.Namespace,
			Labels:      service.Spec.Labels,
			Annotations: service.Spec.Annotations,
		},
		Spec: corev1.ServiceSpec{
			Ports: ports.ToServicePorts(service.Spec.Ports),
			Type:  corev1.ServiceTypeClusterIP,
			Selector: typed.Concat(labels.ManagedByApp(service.Spec.AppNamespace, service.Spec.AppName),
				service.Spec.ContainerLabels),
		},
	}
	result = append(result, newService)
	return
}

func toContainerService(service *v1.ServiceInstance) (result []kclient.Object) {
	newService := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        service.Name,
			Namespace:   service.Namespace,
			Labels:      service.Spec.Labels,
			Annotations: service.Spec.Annotations,
		},
		Spec: corev1.ServiceSpec{
			Ports: ports.ToServicePorts(service.Spec.Ports),
			Type:  corev1.ServiceTypeClusterIP,
			Selector: labels.ManagedByApp(service.Spec.AppNamespace,
				service.Spec.AppName, labels.AcornContainerName, service.Spec.Container),
		},
	}
	result = append(result, newService)
	return
}

func toAddressService(service *v1.ServiceInstance) (result []kclient.Object) {
	newService := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        service.Name,
			Namespace:   service.Namespace,
			Labels:      service.Spec.Labels,
			Annotations: service.Spec.Annotations,
		},
		Spec: corev1.ServiceSpec{
			Ports: ports.ToServicePorts(service.Spec.Ports),
		},
	}
	ipAddr := net.ParseIP(service.Spec.Address)
	if ipAddr == nil {
		newService.Spec.Type = corev1.ServiceTypeExternalName
		newService.Spec.ExternalName = service.Spec.Address
	} else {
		newService.Spec.Type = corev1.ServiceTypeClusterIP
		result = append(result, &corev1.Endpoints{
			ObjectMeta: metav1.ObjectMeta{
				Name:        newService.Name,
				Namespace:   newService.Namespace,
				Labels:      newService.Labels,
				Annotations: newService.Annotations,
			},
			Subsets: []corev1.EndpointSubset{
				{
					Addresses: []corev1.EndpointAddress{
						{
							IP: service.Spec.Address,
						},
					},
					Ports: typed.MapSlice(newService.Spec.Ports, func(t corev1.ServicePort) corev1.EndpointPort {
						return corev1.EndpointPort{
							Name:        t.Name,
							Port:        t.Port,
							Protocol:    t.Protocol,
							AppProtocol: t.AppProtocol,
						}
					}),
				},
			},
		})
	}
	result = append(result, newService)
	return
}

func toExternalService(ctx context.Context, c kclient.Client, cfg *apiv1.Config, service *v1.ServiceInstance) (result []kclient.Object, missing []string, err error) {
	svc, err := resolveTargetService(ctx, c, service.Spec.External, service.Spec.AppNamespace)
	if apierrors.IsNotFound(err) {
		missing = append(missing, service.Spec.External)
	} else if err != nil {
		return nil, nil, err
	}

	if svc == nil || len(svc.Spec.Ports) == 0 {
		return nil, missing, nil
	}

	newService := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        service.Name,
			Namespace:   service.Namespace,
			Labels:      service.Spec.Labels,
			Annotations: service.Spec.Annotations,
		},
		Spec: corev1.ServiceSpec{
			Type:         corev1.ServiceTypeExternalName,
			ExternalName: fmt.Sprintf("%s.%s.%s", svc.Name, svc.Namespace, cfg.InternalClusterDomain),
			Ports:        ports.CopyServicePorts(svc.Spec.Ports),
		},
	}
	result = append(result, newService)
	return
}

func toDefaultService(cfg *apiv1.Config, svc *v1.ServiceInstance, service *corev1.Service) kclient.Object {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        svc.Spec.AppName,
			Namespace:   svc.Spec.AppNamespace,
			Labels:      service.Labels,
			Annotations: service.Annotations,
		},
		Spec: corev1.ServiceSpec{
			Type:         corev1.ServiceTypeExternalName,
			ExternalName: fmt.Sprintf("%s.%s.%s", service.Name, service.Namespace, cfg.InternalClusterDomain),
			Ports:        ports.CopyServicePorts(service.Spec.Ports),
		},
	}
}

func ToK8sService(req router.Request, service *v1.ServiceInstance) (result []kclient.Object, missing []string, err error) {
	cfg, err := config.Get(req.Ctx, req.Client)
	if err != nil {
		return nil, nil, err
	}

	defer func() {
		if err != nil {
			return
		}
		if service.Spec.Default {
			for _, obj := range result {
				if svc, ok := obj.(*corev1.Service); ok {
					result = append(result, toDefaultService(cfg, service, svc))
					return
				}
			}
		}
	}()

	if service.Spec.External != "" {
		return toExternalService(req.Ctx, req.Client, cfg, service)
	} else if service.Spec.Address != "" {
		return toAddressService(service), nil, nil
	} else if service.Spec.Container != "" {
		return toContainerService(service), nil, nil
	} else if len(service.Spec.ContainerLabels) > 0 {
		return toContainerLabelsService(service), nil, nil
	}
	return
}
