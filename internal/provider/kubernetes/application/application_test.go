// Copyright 2020 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package application_test

import (
	"context"
	"fmt"
	stdtesting "testing"
	"time"

	jujuclock "github.com/juju/clock"
	"github.com/juju/clock/testclock"
	"github.com/juju/errors"
	"github.com/juju/tc"
	"go.uber.org/mock/gomock"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	k8sresource "k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/tools/cache"
	"k8s.io/utils/pointer"

	"github.com/juju/juju/caas"
	"github.com/juju/juju/core/constraints"
	"github.com/juju/juju/core/network"
	coreresources "github.com/juju/juju/core/resource"
	"github.com/juju/juju/core/semversion"
	"github.com/juju/juju/core/status"
	"github.com/juju/juju/internal/provider/kubernetes/application"
	"github.com/juju/juju/internal/provider/kubernetes/constants"
	"github.com/juju/juju/internal/provider/kubernetes/resources"
	resourcesmocks "github.com/juju/juju/internal/provider/kubernetes/resources/mocks"
	k8sutils "github.com/juju/juju/internal/provider/kubernetes/utils"
	k8swatcher "github.com/juju/juju/internal/provider/kubernetes/watcher"
	k8swatchertest "github.com/juju/juju/internal/provider/kubernetes/watcher/test"
	"github.com/juju/juju/internal/storage"
	"github.com/juju/juju/internal/testing"
)

type applicationSuite struct {
	testing.BaseSuite
	client *fake.Clientset

	namespace    string
	appName      string
	clock        *testclock.Clock
	k8sWatcherFn k8swatcher.NewK8sWatcherFunc
	watchers     []k8swatcher.KubernetesNotifyWatcher
	applier      *resourcesmocks.MockApplier
}

const defaultAgentVersion = "3.5-beta1"

func TestApplicationSuite(t *stdtesting.T) {
	tc.Run(t, &applicationSuite{})
}

func (s *applicationSuite) SetUpTest(c *tc.C) {
	s.BaseSuite.SetUpTest(c)

	s.namespace = "test"
	s.appName = "gitlab"
	s.client = fake.NewSimpleClientset()
	s.clock = testclock.NewClock(time.Time{})
}

func (s *applicationSuite) TearDownTest(c *tc.C) {
	s.client = nil
	s.clock = nil
	s.watchers = nil
	s.applier = nil

	s.BaseSuite.TearDownTest(c)
}

func (s *applicationSuite) getApp(c *tc.C, deploymentType caas.DeploymentType, mockApplier bool) (application.ApplicationInterfaceForTest, *gomock.Controller) {
	watcherFn := k8swatcher.NewK8sWatcherFunc(func(i cache.SharedIndexInformer, n string, c jujuclock.Clock) (k8swatcher.KubernetesNotifyWatcher, error) {
		if s.k8sWatcherFn == nil {
			return nil, errors.NewNotFound(nil, "undefined k8sWatcherFn for base test")
		}

		w, err := s.k8sWatcherFn(i, n, c)
		if err == nil {
			s.watchers = append(s.watchers, w)
		}
		return w, err
	})

	ctrl := gomock.NewController(c)
	s.applier = resourcesmocks.NewMockApplier(ctrl)

	return application.NewApplicationForTest(
		s.appName, s.namespace, "deadbeef", s.namespace, 2,
		deploymentType,
		s.client,
		watcherFn,
		s.clock,
		func() (string, error) {
			return "appuuid", nil
		},
		func() resources.Applier {
			if mockApplier {
				return s.applier
			}
			return resources.NewApplier()
		},
	), ctrl
}

func (s *applicationSuite) assertEnsure(c *tc.C, app caas.Application, isPrivateImageRepo bool, cons constraints.Value, trust bool, rootless bool, agentVersion string, checkMainResource func()) {
	if agentVersion == "" {
		agentVersion = defaultAgentVersion
	}

	appSecret := corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "gitlab-application-config",
			Namespace: "test",
			Labels: map[string]string{
				"app.kubernetes.io/name":       "gitlab",
				"app.kubernetes.io/managed-by": "juju",
			},
			Annotations: map[string]string{"juju.is/version": agentVersion},
		},
		Data: map[string][]byte{
			"JUJU_K8S_APPLICATION":          []byte("gitlab"),
			"JUJU_K8S_MODEL":                []byte("deadbeef"),
			"JUJU_K8S_APPLICATION_PASSWORD": []byte(""),
			"JUJU_K8S_CONTROLLER_ADDRESSES": []byte(""),
			"JUJU_K8S_CONTROLLER_CA_CERT":   []byte(""),
		},
	}
	appSvc := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "gitlab",
			Namespace: "test",
			Labels: map[string]string{
				"app.kubernetes.io/name":       "gitlab",
				"app.kubernetes.io/managed-by": "juju",
			},
			Annotations: map[string]string{"juju.is/version": agentVersion},
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{"app.kubernetes.io/name": "gitlab"},
			Type:     corev1.ServiceTypeClusterIP,
			Ports: []corev1.ServicePort{{
				Name: "placeholder",
				Port: 65535,
			}},
		},
	}
	pullSecretConfig, _ := k8sutils.CreateDockerConfigJSON("username", "password", "docker.io/library/nginx:latest")
	nginxPullSecret := corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "gitlab-nginx-secret",
			Namespace: "test",
			Labels: map[string]string{
				"app.kubernetes.io/name":       "gitlab",
				"app.kubernetes.io/managed-by": "juju",
			},
			Annotations: map[string]string{"juju.is/version": agentVersion},
		},
		Type: corev1.SecretTypeDockerConfigJson,
		Data: map[string][]byte{
			corev1.DockerConfigJsonKey: pullSecretConfig,
		},
	}
	appSA := corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "gitlab",
			Namespace: "test",
			Labels: map[string]string{
				"app.kubernetes.io/name":       "gitlab",
				"app.kubernetes.io/managed-by": "juju",
			},
			Annotations: map[string]string{"juju.is/version": agentVersion},
		},
		AutomountServiceAccountToken: pointer.BoolPtr(false),
	}
	appRole := rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "gitlab",
			Namespace: "test",
			Labels: map[string]string{
				"app.kubernetes.io/name":       "gitlab",
				"app.kubernetes.io/managed-by": "juju",
			},
			Annotations: map[string]string{"juju.is/version": agentVersion},
		},
	}
	if trust {
		appRole.Rules = []rbacv1.PolicyRule{{
			Verbs:     []string{"*"},
			APIGroups: []string{"*"},
			Resources: []string{"*"},
		}}
	} else {
		appRole.Rules = []rbacv1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"namespaces"},
				Verbs: []string{
					"get",
					"list",
				},
				ResourceNames: []string{s.namespace},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"pods", "services"},
				Verbs: []string{
					"get",
					"list",
					"patch",
				},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"pods/exec"},
				Verbs: []string{
					"create",
				},
			},
		}
	}
	appRoleBinding := rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "gitlab",
			Namespace: "test",
			Labels: map[string]string{
				"app.kubernetes.io/name":       "gitlab",
				"app.kubernetes.io/managed-by": "juju",
			},
			Annotations: map[string]string{"juju.is/version": agentVersion},
		},
		Subjects: []rbacv1.Subject{{
			Kind:      "ServiceAccount",
			Name:      "gitlab",
			Namespace: "test",
		}},
		RoleRef: rbacv1.RoleRef{
			Kind: "Role",
			Name: "gitlab",
		},
	}
	appClusterRole := rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-gitlab",
			Labels: map[string]string{
				"app.kubernetes.io/name":       "gitlab",
				"app.kubernetes.io/managed-by": "juju",
			},
			Annotations: map[string]string{"juju.is/version": agentVersion},
		},
	}
	if trust {
		appClusterRole.Rules = []rbacv1.PolicyRule{{
			Verbs:     []string{"*"},
			APIGroups: []string{"*"},
			Resources: []string{"*"},
		}}
	} else {
		appClusterRole.Rules = []rbacv1.PolicyRule{{
			Verbs:     []string{"get", "list"},
			APIGroups: []string{""},
			Resources: []string{"namespaces"},
		}}
	}
	appClusterRoleBinding := rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-gitlab",
			Labels: map[string]string{
				"app.kubernetes.io/name":       "gitlab",
				"app.kubernetes.io/managed-by": "juju",
			},
			Annotations: map[string]string{"juju.is/version": agentVersion},
		},
		Subjects: []rbacv1.Subject{{
			Kind:      "ServiceAccount",
			Name:      "gitlab",
			Namespace: "test",
		}},
		RoleRef: rbacv1.RoleRef{
			Kind: "ClusterRole",
			Name: "test-gitlab",
		},
	}

	appConfig := caas.ApplicationConfig{
		AgentVersion:         semversion.MustParse(agentVersion),
		IsPrivateImageRepo:   isPrivateImageRepo,
		AgentImagePath:       "operator/image-path:1.1.1",
		CharmBaseImagePath:   "ubuntu@22.04",
		CharmModifiedVersion: 9001,
		Filesystems: []storage.KubernetesFilesystemParams{
			{
				StorageName: "database",
				Size:        100,
				Provider:    "kubernetes",
				Attributes:  map[string]interface{}{"storage-class": "workload-storage"},
				Attachment: &storage.KubernetesFilesystemAttachmentParams{
					Path: "path/to/here",
				},
				ResourceTags: map[string]string{"foo": "bar"},
			},
			// TODO(sidecar): fix here - all filesystems will not be mounted if it's not in `Containers[*].Mounts`
			// {
			// 	StorageName: "logs",
			// 	Size:        200,
			// 	Provider:    "tmpfs",
			// 	Attributes:  map[string]interface{}{"storage-medium": "Memory"},
			// 	Attachment: &storage.KubernetesFilesystemAttachmentParams{
			// 		Path: "path/to/there",
			// 	},
			// },
		},
		Containers: map[string]caas.ContainerConfig{
			"gitlab": {
				Name: "gitlab",
				Image: coreresources.DockerImageDetails{
					RegistryPath: "docker.io/library/gitlab:latest",
				},
				Mounts: []caas.MountConfig{
					{
						StorageName: "database",
						Path:        "path/to/here",
					},
				},
			},
			"nginx": {
				Name: "nginx",
				Image: coreresources.DockerImageDetails{
					RegistryPath: "docker.io/library/nginx:latest",
					ImageRepoDetails: coreresources.ImageRepoDetails{
						BasicAuthConfig: coreresources.BasicAuthConfig{
							Username: "username",
							Password: "password",
						},
					},
				},
				Uid: func() *int {
					if rootless {
						uid := 1234
						return &uid
					}
					return nil
				}(),
				Gid: func() *int {
					if rootless {
						gid := 4321
						return &gid
					}
					return nil
				}(),
			},
		},
		Constraints:  cons,
		InitialScale: 3,
		Trust:        trust,
		CharmUser: func() caas.RunAs {
			if rootless {
				return caas.RunAsNonRoot
			}
			return caas.RunAsDefault
		}(),
	}

	c.Assert(app.Ensure(appConfig), tc.ErrorIsNil)

	secret, err := s.client.CoreV1().Secrets("test").Get(c.Context(), "gitlab-application-config", metav1.GetOptions{})
	c.Assert(err, tc.ErrorIsNil)
	c.Assert(secret, tc.DeepEquals, &appSecret)

	secret, err = s.client.CoreV1().Secrets("test").Get(c.Context(), "gitlab-nginx-secret", metav1.GetOptions{})
	c.Assert(err, tc.ErrorIsNil)
	c.Assert(secret, tc.DeepEquals, &nginxPullSecret)

	svc, err := s.client.CoreV1().Services("test").Get(c.Context(), "gitlab", metav1.GetOptions{})
	c.Assert(err, tc.ErrorIsNil)
	c.Assert(svc, tc.DeepEquals, &appSvc)

	sa, err := s.client.CoreV1().ServiceAccounts(s.namespace).Get(c.Context(), "gitlab", metav1.GetOptions{})
	c.Assert(err, tc.ErrorIsNil)
	c.Assert(sa, tc.DeepEquals, &appSA)

	r, err := s.client.RbacV1().Roles(s.namespace).Get(c.Context(), "gitlab", metav1.GetOptions{})
	c.Assert(err, tc.ErrorIsNil)
	c.Assert(r, tc.DeepEquals, &appRole)

	cr, err := s.client.RbacV1().ClusterRoles().Get(c.Context(), "test-gitlab", metav1.GetOptions{})
	c.Assert(err, tc.ErrorIsNil)
	c.Assert(cr, tc.DeepEquals, &appClusterRole)

	rb, err := s.client.RbacV1().RoleBindings(s.namespace).Get(c.Context(), "gitlab", metav1.GetOptions{})
	c.Assert(err, tc.ErrorIsNil)
	c.Assert(rb, tc.DeepEquals, &appRoleBinding)

	crb, err := s.client.RbacV1().ClusterRoleBindings().Get(c.Context(), "test-gitlab", metav1.GetOptions{})
	c.Assert(err, tc.ErrorIsNil)
	c.Assert(crb, tc.DeepEquals, &appClusterRoleBinding)

	checkMainResource()
}

func (s *applicationSuite) assertDelete(c *tc.C, app caas.Application) {
	err := app.Delete()
	c.Assert(err, tc.ErrorIsNil)

	clusterRoles, err := s.client.RbacV1().ClusterRoles().List(c.Context(), metav1.ListOptions{})
	c.Assert(err, tc.ErrorIsNil)
	c.Assert(clusterRoles.Items, tc.IsNil)

	clusterRoleBinding, err := s.client.RbacV1().ClusterRoleBindings().List(c.Context(), metav1.ListOptions{})
	c.Assert(err, tc.ErrorIsNil)
	c.Assert(clusterRoleBinding.Items, tc.IsNil)

	daemonSets, err := s.client.AppsV1().DaemonSets(s.namespace).List(c.Context(), metav1.ListOptions{})
	c.Assert(err, tc.ErrorIsNil)
	c.Assert(daemonSets.Items, tc.IsNil)

	deployments, err := s.client.AppsV1().Deployments(s.namespace).List(c.Context(), metav1.ListOptions{})
	c.Assert(err, tc.ErrorIsNil)
	c.Assert(deployments.Items, tc.IsNil)

	roles, err := s.client.RbacV1().Roles(s.namespace).List(c.Context(), metav1.ListOptions{})
	c.Assert(err, tc.ErrorIsNil)
	c.Assert(roles.Items, tc.IsNil)

	roleBindings, err := s.client.RbacV1().RoleBindings(s.namespace).List(c.Context(), metav1.ListOptions{})
	c.Assert(err, tc.ErrorIsNil)
	c.Assert(roleBindings.Items, tc.IsNil)

	secrets, err := s.client.CoreV1().Secrets(s.namespace).List(c.Context(), metav1.ListOptions{})
	c.Assert(err, tc.ErrorIsNil)
	c.Assert(secrets.Items, tc.IsNil)

	services, err := s.client.CoreV1().Services(s.namespace).List(c.Context(), metav1.ListOptions{})
	c.Assert(err, tc.ErrorIsNil)
	c.Assert(services.Items, tc.IsNil)

	serviceAccounts, err := s.client.CoreV1().ServiceAccounts(s.namespace).List(c.Context(), metav1.ListOptions{})
	c.Assert(err, tc.ErrorIsNil)
	c.Assert(serviceAccounts.Items, tc.IsNil)

	statefulSets, err := s.client.AppsV1().StatefulSets(s.namespace).List(c.Context(), metav1.ListOptions{})
	c.Assert(err, tc.ErrorIsNil)
	c.Assert(statefulSets.Items, tc.IsNil)
}

func (s *applicationSuite) TestEnsureStateful(c *tc.C) {
	app, _ := s.getApp(c, caas.DeploymentStateful, false)
	s.assertEnsure(
		c, app, false, constraints.Value{}, true, false, "", func() {
			svc, err := s.client.CoreV1().Services("test").Get(c.Context(), "gitlab-endpoints", metav1.GetOptions{})
			c.Assert(err, tc.ErrorIsNil)
			c.Assert(svc, tc.DeepEquals, &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "gitlab-endpoints",
					Namespace: "test",
					Labels: map[string]string{
						"app.kubernetes.io/name":       "gitlab",
						"app.kubernetes.io/managed-by": "juju",
					},
					Annotations: map[string]string{
						"juju.is/version": "3.5-beta1",
						"service.alpha.kubernetes.io/tolerate-unready-endpoints": "true",
					},
				},
				Spec: corev1.ServiceSpec{
					Selector:                 map[string]string{"app.kubernetes.io/name": "gitlab"},
					Type:                     corev1.ServiceTypeClusterIP,
					ClusterIP:                "None",
					PublishNotReadyAddresses: true,
				},
			})

			ss, err := s.client.AppsV1().StatefulSets("test").Get(c.Context(), "gitlab", metav1.GetOptions{})
			c.Assert(err, tc.ErrorIsNil)
			c.Assert(ss, tc.DeepEquals, &appsv1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "gitlab",
					Namespace: "test",
					Labels: map[string]string{
						"app.kubernetes.io/name":       "gitlab",
						"app.kubernetes.io/managed-by": "juju",
					},
					Annotations: map[string]string{
						"juju.is/version":  "3.5-beta1",
						"app.juju.is/uuid": "appuuid",
					},
				},
				Spec: appsv1.StatefulSetSpec{
					Replicas: pointer.Int32Ptr(3),
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"app.kubernetes.io/name": "gitlab",
						},
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels:      map[string]string{"app.kubernetes.io/name": "gitlab"},
							Annotations: map[string]string{"juju.is/version": "3.5-beta1"},
						},
						Spec: getPodSpec31(),
					},
					VolumeClaimTemplates: []corev1.PersistentVolumeClaim{
						{
							ObjectMeta: metav1.ObjectMeta{
								Name: "gitlab-database-appuuid",
								Labels: map[string]string{
									"storage.juju.is/name":         "database",
									"app.kubernetes.io/managed-by": "juju",
								},
								Annotations: map[string]string{
									"foo":                  "bar",
									"storage.juju.is/name": "database",
								}},
							Spec: corev1.PersistentVolumeClaimSpec{
								StorageClassName: pointer.StringPtr("test-workload-storage"),
								AccessModes:      []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
								Resources: corev1.VolumeResourceRequirements{
									Requests: corev1.ResourceList{
										corev1.ResourceStorage: k8sresource.MustParse("100Mi"),
									},
								},
							},
						},
					},
					PodManagementPolicy: appsv1.ParallelPodManagement,
					ServiceName:         "gitlab-endpoints",
				},
			})
		},
	)
	s.assertDelete(c, app)
}

func (s *applicationSuite) TestEnsureStatefulRootless35(c *tc.C) {
	app, _ := s.getApp(c, caas.DeploymentStateful, false)
	s.assertEnsure(
		c, app, false, constraints.Value{}, true, true, "3.5-beta1", func() {
			svc, err := s.client.CoreV1().Services("test").Get(c.Context(), "gitlab-endpoints", metav1.GetOptions{})
			c.Assert(err, tc.ErrorIsNil)
			c.Assert(svc, tc.DeepEquals, &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "gitlab-endpoints",
					Namespace: "test",
					Labels: map[string]string{
						"app.kubernetes.io/name":       "gitlab",
						"app.kubernetes.io/managed-by": "juju",
					},
					Annotations: map[string]string{
						"juju.is/version": "3.5-beta1",
						"service.alpha.kubernetes.io/tolerate-unready-endpoints": "true",
					},
				},
				Spec: corev1.ServiceSpec{
					Selector:                 map[string]string{"app.kubernetes.io/name": "gitlab"},
					Type:                     corev1.ServiceTypeClusterIP,
					ClusterIP:                "None",
					PublishNotReadyAddresses: true,
				},
			})

			podSpec := getPodSpec35()
			podSpec.SecurityContext = &corev1.PodSecurityContext{
				FSGroup:            int64Ptr(170),
				SupplementalGroups: []int64{170},
			}
			podSpec.InitContainers[0].SecurityContext = &corev1.SecurityContext{
				RunAsUser:  int64Ptr(170),
				RunAsGroup: int64Ptr(170),
			}
			podSpec.Containers[0].SecurityContext = &corev1.SecurityContext{
				RunAsUser:  int64Ptr(170),
				RunAsGroup: int64Ptr(170),
			}
			podSpec.Containers[1].SecurityContext = &corev1.SecurityContext{}
			podSpec.Containers[2].SecurityContext = &corev1.SecurityContext{
				RunAsUser:  int64Ptr(1234),
				RunAsGroup: int64Ptr(4321),
			}

			ss, err := s.client.AppsV1().StatefulSets("test").Get(c.Context(), "gitlab", metav1.GetOptions{})
			c.Assert(err, tc.ErrorIsNil)
			c.Assert(ss, tc.DeepEquals, &appsv1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "gitlab",
					Namespace: "test",
					Labels: map[string]string{
						"app.kubernetes.io/name":       "gitlab",
						"app.kubernetes.io/managed-by": "juju",
					},
					Annotations: map[string]string{
						"juju.is/version":  "3.5-beta1",
						"app.juju.is/uuid": "appuuid",
					},
				},
				Spec: appsv1.StatefulSetSpec{
					Replicas: pointer.Int32Ptr(3),
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"app.kubernetes.io/name": "gitlab",
						},
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels:      map[string]string{"app.kubernetes.io/name": "gitlab"},
							Annotations: map[string]string{"juju.is/version": "3.5-beta1"},
						},
						Spec: podSpec,
					},
					VolumeClaimTemplates: []corev1.PersistentVolumeClaim{
						{
							ObjectMeta: metav1.ObjectMeta{
								Name: "gitlab-database-appuuid",
								Labels: map[string]string{
									"storage.juju.is/name":         "database",
									"app.kubernetes.io/managed-by": "juju",
								},
								Annotations: map[string]string{
									"foo":                  "bar",
									"storage.juju.is/name": "database",
								}},
							Spec: corev1.PersistentVolumeClaimSpec{
								StorageClassName: pointer.StringPtr("test-workload-storage"),
								AccessModes:      []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
								Resources: corev1.VolumeResourceRequirements{
									Requests: corev1.ResourceList{
										corev1.ResourceStorage: k8sresource.MustParse("100Mi"),
									},
								},
							},
						},
					},
					PodManagementPolicy: appsv1.ParallelPodManagement,
					ServiceName:         "gitlab-endpoints",
				},
			})
		},
	)
	s.assertDelete(c, app)
}

func (s *applicationSuite) TestEnsureStatefulRootless(c *tc.C) {
	app, _ := s.getApp(c, caas.DeploymentStateful, false)
	s.assertEnsure(
		c, app, false, constraints.Value{}, true, true, "3.6-beta3", func() {
			svc, err := s.client.CoreV1().Services("test").Get(context.TODO(), "gitlab-endpoints", metav1.GetOptions{})
			c.Assert(err, tc.ErrorIsNil)
			c.Assert(svc, tc.DeepEquals, &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "gitlab-endpoints",
					Namespace: "test",
					Labels: map[string]string{
						"app.kubernetes.io/name":       "gitlab",
						"app.kubernetes.io/managed-by": "juju",
					},
					Annotations: map[string]string{
						"juju.is/version": "3.6-beta3",
						"service.alpha.kubernetes.io/tolerate-unready-endpoints": "true",
					},
				},
				Spec: corev1.ServiceSpec{
					Selector:                 map[string]string{"app.kubernetes.io/name": "gitlab"},
					Type:                     corev1.ServiceTypeClusterIP,
					ClusterIP:                "None",
					PublishNotReadyAddresses: true,
				},
			})

			podSpec := getPodSpec36()
			ss, err := s.client.AppsV1().StatefulSets("test").Get(context.TODO(), "gitlab", metav1.GetOptions{})
			c.Assert(err, tc.ErrorIsNil)
			c.Assert(ss, tc.DeepEquals, &appsv1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "gitlab",
					Namespace: "test",
					Labels: map[string]string{
						"app.kubernetes.io/name":       "gitlab",
						"app.kubernetes.io/managed-by": "juju",
					},
					Annotations: map[string]string{
						"juju.is/version":  "3.6-beta3",
						"app.juju.is/uuid": "appuuid",
					},
				},
				Spec: appsv1.StatefulSetSpec{
					Replicas: pointer.Int32Ptr(3),
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"app.kubernetes.io/name": "gitlab",
						},
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels:      map[string]string{"app.kubernetes.io/name": "gitlab"},
							Annotations: map[string]string{"juju.is/version": "3.6-beta3"},
						},
						Spec: podSpec,
					},
					VolumeClaimTemplates: []corev1.PersistentVolumeClaim{
						{
							ObjectMeta: metav1.ObjectMeta{
								Name: "gitlab-database-appuuid",
								Labels: map[string]string{
									"storage.juju.is/name":         "database",
									"app.kubernetes.io/managed-by": "juju",
								},
								Annotations: map[string]string{
									"foo":                  "bar",
									"storage.juju.is/name": "database",
								}},
							Spec: corev1.PersistentVolumeClaimSpec{
								StorageClassName: pointer.StringPtr("test-workload-storage"),
								AccessModes:      []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
								Resources: corev1.VolumeResourceRequirements{
									Requests: corev1.ResourceList{
										corev1.ResourceStorage: k8sresource.MustParse("100Mi"),
									},
								},
							},
						},
					},
					PodManagementPolicy: appsv1.ParallelPodManagement,
					ServiceName:         "gitlab-endpoints",
				},
			})
		},
	)
	s.assertDelete(c, app)
}

func (s *applicationSuite) TestEnsureTrusted(c *tc.C) {
	app, _ := s.getApp(c, caas.DeploymentStateful, false)
	s.assertEnsure(
		c, app, false, constraints.Value{}, true, false, "", func() {},
	)
	s.assertDelete(c, app)
}

func (s *applicationSuite) TestEnsureUntrusted(c *tc.C) {
	app, _ := s.getApp(c, caas.DeploymentStateful, false)
	s.assertEnsure(
		c, app, false, constraints.Value{}, false, false, "", func() {},
	)
	s.assertDelete(c, app)
}

func (s *applicationSuite) TestEnsureStatefulPrivateImageRepo(c *tc.C) {
	app, _ := s.getApp(c, caas.DeploymentStateful, false)

	podSpec := getPodSpec31()
	podSpec.ImagePullSecrets = append(
		[]corev1.LocalObjectReference{
			{Name: constants.CAASImageRepoSecretName},
		},
		podSpec.ImagePullSecrets...,
	)
	s.assertEnsure(
		c, app, true, constraints.Value{}, true, false, "", func() {
			svc, err := s.client.CoreV1().Services("test").Get(c.Context(), "gitlab-endpoints", metav1.GetOptions{})
			c.Assert(err, tc.ErrorIsNil)
			c.Assert(svc, tc.DeepEquals, &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "gitlab-endpoints",
					Namespace: "test",
					Labels: map[string]string{
						"app.kubernetes.io/name":       "gitlab",
						"app.kubernetes.io/managed-by": "juju",
					},
					Annotations: map[string]string{
						"juju.is/version": "3.5-beta1",
						"service.alpha.kubernetes.io/tolerate-unready-endpoints": "true",
					},
				},
				Spec: corev1.ServiceSpec{
					Selector:                 map[string]string{"app.kubernetes.io/name": "gitlab"},
					Type:                     corev1.ServiceTypeClusterIP,
					ClusterIP:                "None",
					PublishNotReadyAddresses: true,
				},
			})

			ss, err := s.client.AppsV1().StatefulSets("test").Get(c.Context(), "gitlab", metav1.GetOptions{})
			c.Assert(err, tc.ErrorIsNil)
			c.Assert(ss, tc.DeepEquals, &appsv1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "gitlab",
					Namespace: "test",
					Labels: map[string]string{
						"app.kubernetes.io/name":       "gitlab",
						"app.kubernetes.io/managed-by": "juju",
					},
					Annotations: map[string]string{
						"juju.is/version":  "3.5-beta1",
						"app.juju.is/uuid": "appuuid",
					},
				},
				Spec: appsv1.StatefulSetSpec{
					Replicas: pointer.Int32Ptr(3),
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"app.kubernetes.io/name": "gitlab",
						},
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels:      map[string]string{"app.kubernetes.io/name": "gitlab"},
							Annotations: map[string]string{"juju.is/version": "3.5-beta1"},
						},
						Spec: podSpec,
					},
					VolumeClaimTemplates: []corev1.PersistentVolumeClaim{
						{
							ObjectMeta: metav1.ObjectMeta{
								Name: "gitlab-database-appuuid",
								Labels: map[string]string{
									"storage.juju.is/name":         "database",
									"app.kubernetes.io/managed-by": "juju",
								},
								Annotations: map[string]string{
									"foo":                  "bar",
									"storage.juju.is/name": "database",
								}},
							Spec: corev1.PersistentVolumeClaimSpec{
								StorageClassName: pointer.StringPtr("test-workload-storage"),
								AccessModes:      []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
								Resources: corev1.VolumeResourceRequirements{
									Requests: corev1.ResourceList{
										corev1.ResourceStorage: k8sresource.MustParse("100Mi"),
									},
								},
							},
						},
					},
					PodManagementPolicy: appsv1.ParallelPodManagement,
					ServiceName:         "gitlab-endpoints",
				},
			})
		},
	)
	s.assertDelete(c, app)
}

func (s *applicationSuite) TestEnsureStateless(c *tc.C) {
	app, _ := s.getApp(c, caas.DeploymentStateless, false)
	s.assertEnsure(
		c, app, false, constraints.Value{}, true, false, "", func() {
			ss, err := s.client.AppsV1().Deployments("test").Get(c.Context(), "gitlab", metav1.GetOptions{})
			c.Assert(err, tc.ErrorIsNil)

			pvc, err := s.client.CoreV1().PersistentVolumeClaims("test").Get(c.Context(), "gitlab-database-appuuid", metav1.GetOptions{})
			c.Assert(err, tc.ErrorIsNil)
			c.Assert(pvc, tc.DeepEquals, &corev1.PersistentVolumeClaim{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "gitlab-database-appuuid",
					Namespace: "test",
					Labels: map[string]string{
						"storage.juju.is/name":         "database",
						"app.kubernetes.io/managed-by": "juju",
					},
					Annotations: map[string]string{
						"foo":                  "bar",
						"storage.juju.is/name": "database",
					},
				},
				Spec: corev1.PersistentVolumeClaimSpec{
					StorageClassName: pointer.StringPtr("test-workload-storage"),
					Resources: corev1.VolumeResourceRequirements{
						Requests: corev1.ResourceList{
							corev1.ResourceStorage: k8sresource.MustParse("100Mi"),
						},
					},
					AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
				},
			})

			podSpec := getPodSpec31()
			podSpec.Volumes = append(podSpec.Volumes, corev1.Volume{
				Name: "gitlab-database-appuuid",
				VolumeSource: corev1.VolumeSource{
					PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
						ClaimName: "gitlab-database-appuuid",
					}},
			})
			c.Assert(ss, tc.DeepEquals, &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "gitlab",
					Namespace: "test",
					Labels: map[string]string{
						"app.kubernetes.io/name":       "gitlab",
						"app.kubernetes.io/managed-by": "juju",
					},
					Annotations: map[string]string{
						"juju.is/version":  "3.5-beta1",
						"app.juju.is/uuid": "appuuid",
					},
				},
				Spec: appsv1.DeploymentSpec{
					Replicas: pointer.Int32Ptr(3),
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{"app.kubernetes.io/name": "gitlab"},
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels:      map[string]string{"app.kubernetes.io/name": "gitlab"},
							Annotations: map[string]string{"juju.is/version": "3.5-beta1"},
						},
						Spec: podSpec,
					},
				},
			})
		},
	)
	s.assertDelete(c, app)
}

func (s *applicationSuite) TestEnsureDaemon(c *tc.C) {
	app, _ := s.getApp(c, caas.DeploymentDaemon, false)
	s.assertEnsure(
		c, app, false, constraints.Value{}, true, false, "", func() {
			ss, err := s.client.AppsV1().DaemonSets("test").Get(c.Context(), "gitlab", metav1.GetOptions{})
			c.Assert(err, tc.ErrorIsNil)

			pvc, err := s.client.CoreV1().PersistentVolumeClaims("test").Get(c.Context(), "gitlab-database-appuuid", metav1.GetOptions{})
			c.Assert(err, tc.ErrorIsNil)
			c.Assert(pvc, tc.DeepEquals, &corev1.PersistentVolumeClaim{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "gitlab-database-appuuid",
					Namespace: "test",
					Labels: map[string]string{
						"storage.juju.is/name":         "database",
						"app.kubernetes.io/managed-by": "juju",
					},
					Annotations: map[string]string{
						"foo":                  "bar",
						"storage.juju.is/name": "database",
					},
				},
				Spec: corev1.PersistentVolumeClaimSpec{
					StorageClassName: pointer.StringPtr("test-workload-storage"),
					Resources: corev1.VolumeResourceRequirements{
						Requests: corev1.ResourceList{
							corev1.ResourceStorage: k8sresource.MustParse("100Mi"),
						},
					},
					AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
				},
			})

			podSpec := getPodSpec31()
			podSpec.Volumes = append(podSpec.Volumes, corev1.Volume{
				Name: "gitlab-database-appuuid",
				VolumeSource: corev1.VolumeSource{
					PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
						ClaimName: "gitlab-database-appuuid",
					}},
			})
			c.Assert(ss, tc.DeepEquals, &appsv1.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "gitlab",
					Namespace: "test",
					Labels: map[string]string{
						"app.kubernetes.io/name":       "gitlab",
						"app.kubernetes.io/managed-by": "juju",
					},
					Annotations: map[string]string{
						"juju.is/version":  "3.5-beta1",
						"app.juju.is/uuid": "appuuid",
					},
				},
				Spec: appsv1.DaemonSetSpec{
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{"app.kubernetes.io/name": "gitlab"},
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels:      map[string]string{"app.kubernetes.io/name": "gitlab"},
							Annotations: map[string]string{"juju.is/version": "3.5-beta1"},
						},
						Spec: podSpec,
					},
				},
			})
		},
	)
	s.assertDelete(c, app)
}

func (s *applicationSuite) TestExistsNotSupported(c *tc.C) {
	app, _ := s.getApp(c, "notsupported", false)
	_, err := app.Exists()
	c.Assert(err, tc.ErrorMatches, `unknown deployment type not supported`)
}

func (s *applicationSuite) TestExistsDeployment(c *tc.C) {
	now := metav1.Now()

	app, _ := s.getApp(c, caas.DeploymentStateless, false)
	// Deployment does not exists.
	result, err := app.Exists()
	c.Assert(err, tc.ErrorIsNil)
	c.Assert(result, tc.DeepEquals, caas.DeploymentState{})

	// ensure a terminating Deployment.
	dr := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "gitlab",
			Namespace: "test",
			Labels: map[string]string{
				"app.kubernetes.io/name":       "gitlab",
				"app.kubernetes.io/managed-by": "juju",
			},
			Annotations: map[string]string{"juju.is/version": "2.9.37"},
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app.kubernetes.io/name": "gitlab"},
			},
		},
	}
	dr.SetDeletionTimestamp(&now)
	_, err = s.client.AppsV1().Deployments("test").Create(c.Context(),
		dr, metav1.CreateOptions{})
	c.Assert(err, tc.ErrorIsNil)

	// Deployment exists and is terminating.
	result, err = app.Exists()
	c.Assert(err, tc.ErrorIsNil)
	c.Assert(result, tc.DeepEquals, caas.DeploymentState{
		Exists: true, Terminating: true,
	})
}

func (s *applicationSuite) TestExistsStatefulSet(c *tc.C) {
	now := metav1.Now()

	app, _ := s.getApp(c, caas.DeploymentStateful, false)
	// Statefulset does not exists.
	result, err := app.Exists()
	c.Assert(err, tc.ErrorIsNil)
	c.Assert(result, tc.DeepEquals, caas.DeploymentState{})

	// ensure a terminating Statefulset.
	sr := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "gitlab",
			Namespace: "test",
			Labels: map[string]string{
				"app.kubernetes.io/name":       "gitlab",
				"app.kubernetes.io/managed-by": "juju",
			},
			Annotations: map[string]string{"juju.is/version": "2.9.37"},
		},
		Spec: appsv1.StatefulSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app.kubernetes.io/name": "gitlab"},
			},
		},
	}
	sr.SetDeletionTimestamp(&now)
	_, err = s.client.AppsV1().StatefulSets("test").Create(c.Context(),
		sr, metav1.CreateOptions{})
	c.Assert(err, tc.ErrorIsNil)

	// Statefulset exists and is terminating.
	result, err = app.Exists()
	c.Assert(err, tc.ErrorIsNil)
	c.Assert(result, tc.DeepEquals, caas.DeploymentState{
		Exists: true, Terminating: true,
	})

}

func (s *applicationSuite) TestExistsDaemonSet(c *tc.C) {
	now := metav1.Now()

	app, _ := s.getApp(c, caas.DeploymentDaemon, false)
	// Daemonset does not exists.
	result, err := app.Exists()
	c.Assert(err, tc.ErrorIsNil)
	c.Assert(result, tc.DeepEquals, caas.DeploymentState{})

	// ensure a terminating Daemonset.
	dmr := &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "gitlab",
			Namespace: "test",
			Labels: map[string]string{
				"app.kubernetes.io/name":       "gitlab",
				"app.kubernetes.io/managed-by": "juju",
			},
			Annotations: map[string]string{"juju.is/version": "2.9.37"},
		},
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app.kubernetes.io/name": "gitlab"},
			},
		},
	}
	dmr.SetDeletionTimestamp(&now)
	_, err = s.client.AppsV1().DaemonSets("test").Create(c.Context(),
		dmr, metav1.CreateOptions{})
	c.Assert(err, tc.ErrorIsNil)

	// Daemonset exists and is terminating.
	result, err = app.Exists()
	c.Assert(err, tc.ErrorIsNil)
	c.Assert(result, tc.DeepEquals, caas.DeploymentState{
		Exists: true, Terminating: true,
	})
}

// Test upgrades are performed by ensure. Regression bug for lp1997253
func (s *applicationSuite) TestUpgradeStateful(c *tc.C) {
	app, _ := s.getApp(c, caas.DeploymentStateful, false)
	s.assertEnsure(c, app, false, constraints.Value{}, true, false, "2.9.34", func() {
		ss, err := s.client.AppsV1().StatefulSets("test").Get(c.Context(), "gitlab", metav1.GetOptions{})
		c.Assert(err, tc.ErrorIsNil)

		c.Assert(len(ss.Spec.Template.Spec.InitContainers), tc.Equals, 1)
		c.Assert(ss.Spec.Template.Spec.InitContainers[0].Args, tc.DeepEquals, []string{
			"init",
			"--data-dir", "/var/lib/juju",
			"--bin-dir", "/charm/bin",
		})
	})

	s.assertEnsure(c, app, false, constraints.Value{}, true, false, "2.9.37", func() {
		ss, err := s.client.AppsV1().StatefulSets("test").Get(c.Context(), "gitlab", metav1.GetOptions{})
		c.Assert(err, tc.ErrorIsNil)

		c.Assert(len(ss.Spec.Template.Spec.InitContainers), tc.Equals, 1)
		c.Assert(ss.Spec.Template.Spec.InitContainers[0].Args, tc.DeepEquals, []string{
			"init",
			"--containeragent-pebble-dir", "/containeragent/pebble",
			"--charm-modified-version", "9001",
			"--data-dir", "/var/lib/juju",
			"--bin-dir", "/charm/bin",
		})
	})

	s.assertEnsure(c, app, false, constraints.Value{}, true, false, "3.5-beta1.1", func() {
		ss, err := s.client.AppsV1().StatefulSets("test").Get(c.Context(), "gitlab", metav1.GetOptions{})
		c.Assert(err, tc.ErrorIsNil)

		c.Assert(len(ss.Spec.Template.Spec.InitContainers), tc.Equals, 1)
		c.Assert(ss.Spec.Template.Spec.InitContainers[0].Args, tc.DeepEquals, []string{
			"init",
			"--containeragent-pebble-dir", "/containeragent/pebble",
			"--charm-modified-version", "9001",
			"--data-dir", "/var/lib/juju",
			"--bin-dir", "/charm/bin",
			"--profile-dir", "/containeragent/etc/profile.d",
		})
	})
}

func (s *applicationSuite) TestDeleteStateful(c *tc.C) {
	app, ctrl := s.getApp(c, caas.DeploymentStateful, true)
	defer ctrl.Finish()

	gomock.InOrder(
		s.applier.EXPECT().Delete(resources.NewStatefulSet("gitlab", "test", nil)),
		s.applier.EXPECT().Delete(resources.NewService("gitlab-endpoints", "test", nil)),
		s.applier.EXPECT().Delete(resources.NewService("gitlab", "test", nil)),
		s.applier.EXPECT().Delete(resources.NewSecret("gitlab-application-config", "test", nil)),
		s.applier.EXPECT().Delete(resources.NewRoleBinding("gitlab", "test", nil)),
		s.applier.EXPECT().Delete(resources.NewRole("gitlab", "test", nil)),
		s.applier.EXPECT().Delete(resources.NewClusterRoleBinding("test-gitlab", nil)),
		s.applier.EXPECT().Delete(resources.NewClusterRole("test-gitlab", nil)),
		s.applier.EXPECT().Delete(resources.NewServiceAccount("gitlab", "test", nil)),
		s.applier.EXPECT().Run(gomock.Any(), s.client, false).Return(nil),
	)
	c.Assert(app.Delete(), tc.ErrorIsNil)
}

func (s *applicationSuite) TestDeleteStateless(c *tc.C) {
	app, ctrl := s.getApp(c, caas.DeploymentStateless, true)
	defer ctrl.Finish()

	gomock.InOrder(
		s.applier.EXPECT().Delete(resources.NewDeployment("gitlab", "test", nil)),
		s.applier.EXPECT().Delete(resources.NewService("gitlab", "test", nil)),
		s.applier.EXPECT().Delete(resources.NewSecret("gitlab-application-config", "test", nil)),
		s.applier.EXPECT().Delete(resources.NewRoleBinding("gitlab", "test", nil)),
		s.applier.EXPECT().Delete(resources.NewRole("gitlab", "test", nil)),
		s.applier.EXPECT().Delete(resources.NewClusterRoleBinding("test-gitlab", nil)),
		s.applier.EXPECT().Delete(resources.NewClusterRole("test-gitlab", nil)),
		s.applier.EXPECT().Delete(resources.NewServiceAccount("gitlab", "test", nil)),
		s.applier.EXPECT().Run(gomock.Any(), s.client, false).Return(nil),
	)
	c.Assert(app.Delete(), tc.ErrorIsNil)
}

func (s *applicationSuite) TestDeleteDaemon(c *tc.C) {
	app, ctrl := s.getApp(c, caas.DeploymentDaemon, true)
	defer ctrl.Finish()

	gomock.InOrder(
		s.applier.EXPECT().Delete(resources.NewDaemonSet("gitlab", "test", nil)),
		s.applier.EXPECT().Delete(resources.NewService("gitlab", "test", nil)),
		s.applier.EXPECT().Delete(resources.NewSecret("gitlab-application-config", "test", nil)),
		s.applier.EXPECT().Delete(resources.NewRoleBinding("gitlab", "test", nil)),
		s.applier.EXPECT().Delete(resources.NewRole("gitlab", "test", nil)),
		s.applier.EXPECT().Delete(resources.NewClusterRoleBinding("test-gitlab", nil)),
		s.applier.EXPECT().Delete(resources.NewClusterRole("test-gitlab", nil)),
		s.applier.EXPECT().Delete(resources.NewServiceAccount("gitlab", "test", nil)),
		s.applier.EXPECT().Run(gomock.Any(), s.client, false).Return(nil),
	)
	c.Assert(app.Delete(), tc.ErrorIsNil)
}

func (s *applicationSuite) TestWatchNotsupported(c *tc.C) {
	app, ctrl := s.getApp(c, "notsupported", true)
	defer ctrl.Finish()

	s.k8sWatcherFn = func(_ cache.SharedIndexInformer, _ string, _ jujuclock.Clock) (k8swatcher.KubernetesNotifyWatcher, error) {
		w, _ := k8swatchertest.NewKubernetesTestWatcher()
		return w, nil
	}

	_, err := app.Watch(c.Context())
	c.Assert(err, tc.ErrorMatches, `unknown deployment type not supported`)
}

func (s *applicationSuite) TestWatch(c *tc.C) {
	app, ctrl := s.getApp(c, caas.DeploymentDaemon, true)
	defer ctrl.Finish()

	s.k8sWatcherFn = func(_ cache.SharedIndexInformer, _ string, _ jujuclock.Clock) (k8swatcher.KubernetesNotifyWatcher, error) {
		w, _ := k8swatchertest.NewKubernetesTestWatcher()
		return w, nil
	}

	w, err := app.Watch(c.Context())
	c.Assert(err, tc.ErrorIsNil)

	select {
	case _, ok := <-w.Changes():
		c.Assert(ok, tc.IsTrue)
	case <-time.After(testing.LongWait):
		c.Fatal("timed out waiting for event")
	}
}

func (s *applicationSuite) TestWatchReplicas(c *tc.C) {
	app, ctrl := s.getApp(c, caas.DeploymentDaemon, true)
	defer ctrl.Finish()

	s.k8sWatcherFn = func(_ cache.SharedIndexInformer, _ string, _ jujuclock.Clock) (k8swatcher.KubernetesNotifyWatcher, error) {
		w, _ := k8swatchertest.NewKubernetesTestWatcher()
		return w, nil
	}

	w, err := app.WatchReplicas()
	c.Assert(err, tc.ErrorIsNil)

	select {
	case _, ok := <-w.Changes():
		c.Assert(ok, tc.IsTrue)
	case <-time.After(testing.LongWait):
		c.Fatal("timed out waiting for event")
	}
}

func (s *applicationSuite) TestStateNotSupported(c *tc.C) {
	app, _ := s.getApp(c, "notsupported", false)
	_, err := app.State()
	c.Assert(err, tc.ErrorMatches, `unknown deployment type not supported`)
}

func (s *applicationSuite) assertState(c *tc.C, deploymentType caas.DeploymentType, createMainResource func() int) {
	app, ctrl := s.getApp(c, deploymentType, false)
	defer ctrl.Finish()

	desiredReplicas := createMainResource()

	pod1 := &corev1.Pod{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:        "pod1",
			Namespace:   "test",
			Labels:      map[string]string{"app.kubernetes.io/name": "gitlab"},
			Annotations: map[string]string{"juju.is/version": "2.9.37"},
		},
	}
	_, err := s.client.CoreV1().Pods("test").Create(c.Context(),
		pod1, metav1.CreateOptions{})
	c.Assert(err, tc.ErrorIsNil)

	pod2 := &corev1.Pod{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:        "pod2",
			Namespace:   "test",
			Labels:      map[string]string{"app.kubernetes.io/name": "gitlab"},
			Annotations: map[string]string{"juju.is/version": "2.9.37"},
		},
	}
	_, err = s.client.CoreV1().Pods("test").Create(c.Context(),
		pod2, metav1.CreateOptions{})
	c.Assert(err, tc.ErrorIsNil)

	appState, err := app.State()
	c.Assert(err, tc.ErrorIsNil)
	c.Assert(appState, tc.DeepEquals, caas.ApplicationState{
		DesiredReplicas: desiredReplicas,
		Replicas:        []string{"pod1", "pod2"},
	})
}

func (s *applicationSuite) TestStateStateful(c *tc.C) {
	s.assertState(c, caas.DeploymentStateful, func() int {
		desiredReplicas := 10

		dmr := &appsv1.StatefulSet{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "gitlab",
				Namespace: "test",
				Labels: map[string]string{
					"app.kubernetes.io/name":       "gitlab",
					"app.kubernetes.io/managed-by": "juju",
				},
				Annotations: map[string]string{"juju.is/version": "2.9.37"},
			},
			Spec: appsv1.StatefulSetSpec{
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{"app.kubernetes.io/name": "gitlab"},
				},
				Replicas: pointer.Int32Ptr(int32(desiredReplicas)),
			},
		}
		_, err := s.client.AppsV1().StatefulSets("test").Create(c.Context(),
			dmr, metav1.CreateOptions{})
		c.Assert(err, tc.ErrorIsNil)
		return desiredReplicas
	})
}

func (s *applicationSuite) TestStateStateless(c *tc.C) {
	s.assertState(c, caas.DeploymentStateless, func() int {
		desiredReplicas := 10

		dmr := &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "gitlab",
				Namespace: "test",
				Labels: map[string]string{
					"app.kubernetes.io/name":       "gitlab",
					"app.kubernetes.io/managed-by": "juju",
				},
				Annotations: map[string]string{"juju.is/version": "2.9.37"},
			},
			Spec: appsv1.DeploymentSpec{
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{"app.kubernetes.io/name": "gitlab"},
				},
				Replicas: pointer.Int32Ptr(int32(desiredReplicas)),
			},
		}
		_, err := s.client.AppsV1().Deployments("test").Create(c.Context(),
			dmr, metav1.CreateOptions{})
		c.Assert(err, tc.ErrorIsNil)
		return desiredReplicas
	})
}

func (s *applicationSuite) TestStateDaemon(c *tc.C) {
	s.assertState(c, caas.DeploymentDaemon, func() int {
		desiredReplicas := 10

		dmr := &appsv1.DaemonSet{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "gitlab",
				Namespace: "test",
				Labels: map[string]string{
					"app.kubernetes.io/name":       "gitlab",
					"app.kubernetes.io/managed-by": "juju",
				},
				Annotations: map[string]string{"juju.is/version": "2.9.37"},
			},
			Spec: appsv1.DaemonSetSpec{
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{"app.kubernetes.io/name": "gitlab"},
				},
			},
			Status: appsv1.DaemonSetStatus{
				DesiredNumberScheduled: int32(desiredReplicas),
			},
		}
		_, err := s.client.AppsV1().DaemonSets("test").Create(c.Context(),
			dmr, metav1.CreateOptions{})
		c.Assert(err, tc.ErrorIsNil)
		return desiredReplicas
	})
}

func getDefaultSvc() *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "gitlab",
			Namespace: "test",
			Labels: map[string]string{
				"app.kubernetes.io/name":       "gitlab",
				"app.kubernetes.io/managed-by": "juju",
			},
			Annotations: map[string]string{"juju.is/version": "2.9.37"},
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{"app.kubernetes.io/name": "gitlab"},

			Type: corev1.ServiceTypeClusterIP,
			Ports: []corev1.ServicePort{{
				Name: "placeholder",
				Port: 65535,
			}},
		},
	}
}

func (s *applicationSuite) TestUpdatePortsStatelessUpdateContainerPorts(c *tc.C) {
	app, ctrl := s.getApp(c, caas.DeploymentStateless, true)
	defer ctrl.Finish()

	_, err := s.client.CoreV1().Services("test").Create(c.Context(), getDefaultSvc(), metav1.CreateOptions{})
	c.Assert(err, tc.ErrorIsNil)

	getMainResourceSpec := func() *appsv1.Deployment {
		return &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "gitlab",
				Namespace: "test",
				Labels: map[string]string{
					"app.kubernetes.io/name":       "gitlab",
					"app.kubernetes.io/managed-by": "juju",
				},
				Annotations: map[string]string{
					"juju.is/version":  "2.9.37",
					"app.juju.is/uuid": "appuuid",
				},
			},
			Spec: appsv1.DeploymentSpec{
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{"app.kubernetes.io/name": "gitlab"},
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels:      map[string]string{"app.kubernetes.io/name": "gitlab"},
						Annotations: map[string]string{"juju.is/version": "2.9.37"},
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{{
							Name:            "charm",
							ImagePullPolicy: corev1.PullIfNotPresent,
							Image:           "operator/image-path",
							WorkingDir:      "/var/lib/juju",
							Command:         []string{"/charm/bin/containeragent"},
							Args:            []string{"unit", "--data-dir", "/var/lib/juju", "--append-env", "PATH=$PATH:/charm/bin"},
							Env: []corev1.EnvVar{{
								Name:  "HTTP_PROBE_PORT",
								Value: "3856",
							}},
							LivenessProbe: &corev1.Probe{
								ProbeHandler: corev1.ProbeHandler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/liveness",
										Port: intstr.FromString("3856"),
									},
								},
								InitialDelaySeconds: 30,
								PeriodSeconds:       10,
								SuccessThreshold:    1,
								FailureThreshold:    2,
							},
							ReadinessProbe: &corev1.Probe{
								ProbeHandler: corev1.ProbeHandler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/readiness",
										Port: intstr.FromString("3856"),
									},
								},
								InitialDelaySeconds: 30,
								PeriodSeconds:       10,
								SuccessThreshold:    1,
								FailureThreshold:    2,
							},
							StartupProbe: &corev1.Probe{
								ProbeHandler: corev1.ProbeHandler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/startup",
										Port: intstr.FromString("3856"),
									},
								},
								InitialDelaySeconds: 30,
								PeriodSeconds:       10,
								SuccessThreshold:    1,
								FailureThreshold:    2,
							},
						}, {
							Name:            "gitlab",
							ImagePullPolicy: corev1.PullIfNotPresent,
							Image:           "test-image",
							Command:         []string{"/charm/bin/pebble"},
							Args:            []string{"listen", "--socket", "/charm/container/pebble.sock", "--append-env", "PATH=$PATH:/charm/bin"},
						}},
					},
				},
			},
		}
	}
	_, err = s.client.AppsV1().Deployments("test").Create(c.Context(), getMainResourceSpec(), metav1.CreateOptions{})
	c.Assert(err, tc.ErrorIsNil)

	updatedSvc := getDefaultSvc()
	updatedSvc.Spec.Ports = []corev1.ServicePort{
		{
			Name:       "juju-port1",
			Port:       int32(8080),
			TargetPort: intstr.FromInt(8080),
			Protocol:   corev1.ProtocolTCP,
		},
	}
	updatedSvcResource := resources.NewService("gitlab", "test", updatedSvc)
	replacePortsPatchType := types.MergePatchType
	updatedSvcResource.PatchType = &replacePortsPatchType

	updatedMainResource := getMainResourceSpec()
	updatedMainResource.Spec.Template.Spec.Containers[1].Ports = []corev1.ContainerPort{
		{
			Name:          "juju-port1",
			ContainerPort: int32(8080),
			Protocol:      corev1.ProtocolTCP,
		},
	}
	gomock.InOrder(
		s.applier.EXPECT().Apply(updatedSvcResource),
		s.applier.EXPECT().Apply(resources.NewDeployment("gitlab", "test", updatedMainResource)),
		s.applier.EXPECT().Run(gomock.Any(), s.client, false).Return(nil),
	)
	c.Assert(app.UpdatePorts([]caas.ServicePort{
		{
			Name:       "port1",
			Port:       8080,
			TargetPort: 8080,
			Protocol:   "TCP",
		},
	}, true), tc.ErrorIsNil)
}

func (s *applicationSuite) TestUpdatePortsStatefulUpdateContainerPorts(c *tc.C) {
	app, ctrl := s.getApp(c, caas.DeploymentStateful, true)
	defer ctrl.Finish()

	_, err := s.client.CoreV1().Services("test").Create(c.Context(), getDefaultSvc(), metav1.CreateOptions{})
	c.Assert(err, tc.ErrorIsNil)

	getMainResourceSpec := func() *appsv1.StatefulSet {
		return &appsv1.StatefulSet{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "gitlab",
				Namespace: "test",
				Labels: map[string]string{
					"app.kubernetes.io/name":       "gitlab",
					"app.kubernetes.io/managed-by": "juju",
				},
				Annotations: map[string]string{
					"juju.is/version":  "2.9.37",
					"app.juju.is/uuid": "appuuid",
				},
			},
			Spec: appsv1.StatefulSetSpec{
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{"app.kubernetes.io/name": "gitlab"},
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels:      map[string]string{"app.kubernetes.io/name": "gitlab"},
						Annotations: map[string]string{"juju.is/version": "2.9.37"},
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{{
							Name:            "charm",
							ImagePullPolicy: corev1.PullIfNotPresent,
							Image:           "operator/image-path",
							WorkingDir:      "/var/lib/juju",
							Command:         []string{"/charm/bin/containeragent"},
							Args:            []string{"unit", "--data-dir", "/var/lib/juju", "--append-env", "PATH=$PATH:/charm/bin"},
							Env: []corev1.EnvVar{{
								Name:  "HTTP_PROBE_PORT",
								Value: "3856",
							}},
							LivenessProbe: &corev1.Probe{
								ProbeHandler: corev1.ProbeHandler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/liveness",
										Port: intstr.FromString("3856"),
									},
								},
								InitialDelaySeconds: 30,
								PeriodSeconds:       10,
								SuccessThreshold:    1,
								FailureThreshold:    2,
							},
							ReadinessProbe: &corev1.Probe{
								ProbeHandler: corev1.ProbeHandler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/readiness",
										Port: intstr.FromString("3856"),
									},
								},
								InitialDelaySeconds: 30,
								PeriodSeconds:       10,
								SuccessThreshold:    1,
								FailureThreshold:    2,
							},
							StartupProbe: &corev1.Probe{
								ProbeHandler: corev1.ProbeHandler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/startup",
										Port: intstr.FromString("3856"),
									},
								},
								InitialDelaySeconds: 30,
								PeriodSeconds:       10,
								SuccessThreshold:    1,
								FailureThreshold:    2,
							},
						}, {
							Name:            "gitlab",
							ImagePullPolicy: corev1.PullIfNotPresent,
							Image:           "test-image",
							Command:         []string{"/charm/bin/pebble"},
							Args:            []string{"listen", "--socket", "/charm/container/pebble.sock", "--append-env", "PATH=$PATH:/charm/bin"},
						}},
					},
				},
			},
		}
	}
	_, err = s.client.AppsV1().StatefulSets("test").Create(c.Context(), getMainResourceSpec(), metav1.CreateOptions{})
	c.Assert(err, tc.ErrorIsNil)

	updatedSvc := getDefaultSvc()
	updatedSvc.Spec.Ports = []corev1.ServicePort{
		{
			Name:       "juju-port1",
			Port:       int32(8080),
			TargetPort: intstr.FromInt(8080),
			Protocol:   corev1.ProtocolTCP,
		},
	}
	updatedSvcResource := resources.NewService("gitlab", "test", updatedSvc)
	replacePortsPatchType := types.MergePatchType
	updatedSvcResource.PatchType = &replacePortsPatchType

	updatedMainResource := getMainResourceSpec()
	updatedMainResource.Spec.Template.Spec.Containers[1].Ports = []corev1.ContainerPort{
		{
			Name:          "juju-port1",
			ContainerPort: int32(8080),
			Protocol:      corev1.ProtocolTCP,
		},
	}
	gomock.InOrder(
		s.applier.EXPECT().Apply(updatedSvcResource),
		s.applier.EXPECT().Apply(resources.NewStatefulSet("gitlab", "test", updatedMainResource)),
		s.applier.EXPECT().Run(gomock.Any(), s.client, false).Return(nil),
	)
	c.Assert(app.UpdatePorts([]caas.ServicePort{
		{
			Name:       "port1",
			Port:       8080,
			TargetPort: 8080,
			Protocol:   "TCP",
		},
	}, true), tc.ErrorIsNil)
}

func (s *applicationSuite) TestUpdatePortsDaemonUpdateContainerPorts(c *tc.C) {
	app, ctrl := s.getApp(c, caas.DeploymentDaemon, true)
	defer ctrl.Finish()

	_, err := s.client.CoreV1().Services("test").Create(c.Context(), getDefaultSvc(), metav1.CreateOptions{})
	c.Assert(err, tc.ErrorIsNil)

	getMainResourceSpec := func() *appsv1.DaemonSet {
		return &appsv1.DaemonSet{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "gitlab",
				Namespace: "test",
				Labels: map[string]string{
					"app.kubernetes.io/name":       "gitlab",
					"app.kubernetes.io/managed-by": "juju",
				},
				Annotations: map[string]string{
					"juju.is/version":  "2.9.37",
					"app.juju.is/uuid": "appuuid",
				},
			},
			Spec: appsv1.DaemonSetSpec{
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{"app.kubernetes.io/name": "gitlab"},
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels:      map[string]string{"app.kubernetes.io/name": "gitlab"},
						Annotations: map[string]string{"juju.is/version": "2.9.37"},
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{{
							Name:            "charm",
							ImagePullPolicy: corev1.PullIfNotPresent,
							Image:           "operator/image-path",
							WorkingDir:      "/var/lib/juju",
							Command:         []string{"/charm/bin/containeragent"},
							Args:            []string{"unit", "--data-dir", "/var/lib/juju", "--append-env", "PATH=$PATH:/charm/bin"},
						}, {
							Name:            "gitlab",
							ImagePullPolicy: corev1.PullIfNotPresent,
							Image:           "test-image",
							Command:         []string{"/charm/bin/pebble"},
							Args:            []string{"listen", "--socket", "/charm/container/pebble.sock", "--append-env", "PATH=$PATH:/charm/bin"},
						}},
					},
				},
			},
		}
	}
	_, err = s.client.AppsV1().DaemonSets("test").Create(c.Context(), getMainResourceSpec(), metav1.CreateOptions{})
	c.Assert(err, tc.ErrorIsNil)

	updatedSvc := getDefaultSvc()
	updatedSvc.Spec.Ports = []corev1.ServicePort{
		{
			Name:       "juju-port1",
			Port:       int32(8080),
			TargetPort: intstr.FromInt(8080),
			Protocol:   corev1.ProtocolTCP,
		},
	}
	updatedSvcResource := resources.NewService("gitlab", "test", updatedSvc)
	replacePortsPatchType := types.MergePatchType
	updatedSvcResource.PatchType = &replacePortsPatchType

	updatedMainResource := getMainResourceSpec()
	updatedMainResource.Spec.Template.Spec.Containers[1].Ports = []corev1.ContainerPort{
		{
			Name:          "juju-port1",
			ContainerPort: int32(8080),
			Protocol:      corev1.ProtocolTCP,
		},
	}
	gomock.InOrder(
		s.applier.EXPECT().Apply(updatedSvcResource),
		s.applier.EXPECT().Apply(resources.NewDaemonSet("gitlab", "test", updatedMainResource)),
		s.applier.EXPECT().Run(gomock.Any(), s.client, false).Return(nil),
	)
	c.Assert(app.UpdatePorts([]caas.ServicePort{
		{
			Name:       "port1",
			Port:       8080,
			TargetPort: 8080,
			Protocol:   "TCP",
		},
	}, true), tc.ErrorIsNil)
}

func (s *applicationSuite) TestUpdatePortsInvalidProtocol(c *tc.C) {
	app, ctrl := s.getApp(c, caas.DeploymentStateful, true)
	defer ctrl.Finish()

	_, err := s.client.CoreV1().Services("test").Create(c.Context(), getDefaultSvc(), metav1.CreateOptions{})
	c.Assert(err, tc.ErrorIsNil)

	c.Assert(app.UpdatePorts([]caas.ServicePort{
		{
			Name:       "port1",
			Port:       8080,
			TargetPort: 8080,
			Protocol:   "bad-protocol",
		},
	}, false), tc.ErrorMatches, `protocol "bad-protocol" for service "port1" not valid`)
}

func (s *applicationSuite) TestUpdatePortsWithExistingPorts(c *tc.C) {
	app, ctrl := s.getApp(c, caas.DeploymentStateful, true)
	defer ctrl.Finish()

	existingSvc := getDefaultSvc()
	existingSvc.Spec.Ports = []corev1.ServicePort{
		{
			Name:       "existing-port",
			Port:       int32(3000),
			TargetPort: intstr.FromInt(3000),
			Protocol:   corev1.ProtocolTCP,
		},
	}
	svc, err := s.client.CoreV1().Services("test").Create(c.Context(), existingSvc, metav1.CreateOptions{})
	c.Assert(err, tc.ErrorIsNil)
	c.Assert(svc.Spec.Ports, tc.DeepEquals, existingSvc.Spec.Ports)

	updatedSvc := getDefaultSvc()
	updatedSvc.Spec.Ports = []corev1.ServicePort{
		{
			Name:       "existing-port",
			Port:       int32(3000),
			TargetPort: intstr.FromInt(3000),
			Protocol:   corev1.ProtocolTCP,
		},
		{
			Name:       "juju-port1",
			Port:       int32(8080),
			TargetPort: intstr.FromInt(8080),
			Protocol:   corev1.ProtocolTCP,
		},
		{
			Name:       "juju-port2",
			Port:       int32(8888),
			TargetPort: intstr.FromInt(8888),
			Protocol:   corev1.ProtocolTCP,
		},
	}
	updatedSvcResource := resources.NewService("gitlab", "test", updatedSvc)
	replacePortsPatchType := types.MergePatchType
	updatedSvcResource.PatchType = &replacePortsPatchType

	updatedSvc2nd := getDefaultSvc()
	updatedSvc2nd.Spec.Ports = []corev1.ServicePort{
		{
			Name:       "existing-port",
			Port:       int32(3000),
			TargetPort: intstr.FromInt(3000),
			Protocol:   corev1.ProtocolTCP,
		},
		{
			Name:       "juju-port2",
			Port:       int32(8888),
			TargetPort: intstr.FromInt(8888),
			Protocol:   corev1.ProtocolTCP,
		},
	}
	updatedSvcResource2nd := resources.NewService("gitlab", "test", updatedSvc2nd)
	updatedSvcResource2nd.PatchType = &replacePortsPatchType

	gomock.InOrder(
		s.applier.EXPECT().Apply(updatedSvcResource),
		s.applier.EXPECT().Run(gomock.Any(), s.client, false).Return(nil),

		s.applier.EXPECT().Apply(updatedSvcResource2nd),
		s.applier.EXPECT().Run(gomock.Any(), s.client, false).Return(nil),
	)
	c.Assert(app.UpdatePorts([]caas.ServicePort{
		// Added ports: 8080 and 8888.
		{
			Name:       "port1",
			Port:       8080,
			TargetPort: 8080,
			Protocol:   "TCP",
		},
		{
			Name:       "port2",
			Port:       8888,
			TargetPort: 8888,
			Protocol:   "TCP",
		},
	}, false), tc.ErrorIsNil)

	c.Assert(app.UpdatePorts([]caas.ServicePort{
		// Removed port: 8080.
		{
			Name:       "port2",
			Port:       8888,
			TargetPort: 8888,
			Protocol:   "TCP",
		},
	}, false), tc.ErrorIsNil)
}

func (s *applicationSuite) TestUpdatePortsStateless(c *tc.C) {
	app, ctrl := s.getApp(c, caas.DeploymentStateless, true)
	defer ctrl.Finish()

	_, err := s.client.CoreV1().Services("test").Create(c.Context(), getDefaultSvc(), metav1.CreateOptions{})
	c.Assert(err, tc.ErrorIsNil)

	updatedSvc := getDefaultSvc()
	updatedSvc.Spec.Ports = []corev1.ServicePort{
		{
			Name:       "juju-port1",
			Port:       int32(8080),
			TargetPort: intstr.FromInt(8080),
			Protocol:   corev1.ProtocolTCP,
		},
	}
	updatedSvcResource := resources.NewService("gitlab", "test", updatedSvc)
	replacePortsPatchType := types.MergePatchType
	updatedSvcResource.PatchType = &replacePortsPatchType

	gomock.InOrder(
		s.applier.EXPECT().Apply(updatedSvcResource),
		s.applier.EXPECT().Run(gomock.Any(), s.client, false).Return(nil),
	)
	c.Assert(app.UpdatePorts([]caas.ServicePort{
		{
			Name:       "port1",
			Port:       8080,
			TargetPort: 8080,
			Protocol:   "TCP",
		},
	}, false), tc.ErrorIsNil)
}

func (s *applicationSuite) TestUpdatePortsStateful(c *tc.C) {
	app, ctrl := s.getApp(c, caas.DeploymentStateful, true)
	defer ctrl.Finish()

	_, err := s.client.CoreV1().Services("test").Create(c.Context(), getDefaultSvc(), metav1.CreateOptions{})
	c.Assert(err, tc.ErrorIsNil)

	updatedSvc := getDefaultSvc()
	updatedSvc.Spec.Ports = []corev1.ServicePort{
		{
			Name:       "juju-port1",
			Port:       int32(8080),
			TargetPort: intstr.FromInt(8080),
			Protocol:   corev1.ProtocolTCP,
		},
	}
	updatedSvcResource := resources.NewService("gitlab", "test", updatedSvc)
	replacePortsPatchType := types.MergePatchType
	updatedSvcResource.PatchType = &replacePortsPatchType

	gomock.InOrder(
		s.applier.EXPECT().Apply(updatedSvcResource),
		s.applier.EXPECT().Run(gomock.Any(), s.client, false).Return(nil),
	)
	c.Assert(app.UpdatePorts([]caas.ServicePort{
		{
			Name:       "port1",
			Port:       8080,
			TargetPort: 8080,
			Protocol:   "TCP",
		},
	}, false), tc.ErrorIsNil)
}

func (s *applicationSuite) TestUpdatePortsDaemonUpdate(c *tc.C) {
	app, ctrl := s.getApp(c, caas.DeploymentDaemon, true)
	defer ctrl.Finish()

	_, err := s.client.CoreV1().Services("test").Create(c.Context(), getDefaultSvc(), metav1.CreateOptions{})
	c.Assert(err, tc.ErrorIsNil)

	updatedSvc := getDefaultSvc()
	updatedSvc.Spec.Ports = []corev1.ServicePort{
		{
			Name:       "juju-port1",
			Port:       int32(8080),
			TargetPort: intstr.FromInt(8080),
			Protocol:   corev1.ProtocolTCP,
		},
	}
	updatedSvcResource := resources.NewService("gitlab", "test", updatedSvc)
	replacePortsPatchType := types.MergePatchType
	updatedSvcResource.PatchType = &replacePortsPatchType

	gomock.InOrder(
		s.applier.EXPECT().Apply(updatedSvcResource),
		s.applier.EXPECT().Run(gomock.Any(), s.client, false).Return(nil),
	)
	c.Assert(app.UpdatePorts([]caas.ServicePort{
		{
			Name:       "port1",
			Port:       8080,
			TargetPort: 8080,
			Protocol:   "TCP",
		},
	}, false), tc.ErrorIsNil)
}

func (s *applicationSuite) TestUnits(c *tc.C) {
	app, _ := s.getApp(c, caas.DeploymentStateful, false)

	for i := 0; i < 9; i++ {
		podSpec := getPodSpec31()
		podSpec.Volumes = append(podSpec.Volumes,
			corev1.Volume{
				Name: "gitlab-database-appuuid",
				VolumeSource: corev1.VolumeSource{
					PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
						ClaimName: fmt.Sprintf("gitlab-database-appuuid-gitlab-%d", i),
					},
				},
			},
		)
		// Ensure these volume sources are ignored
		podSpec.Volumes = append(podSpec.Volumes,
			corev1.Volume{
				Name: "vol-secret",
				VolumeSource: corev1.VolumeSource{
					EmptyDir: &corev1.EmptyDirVolumeSource{},
					Secret: &corev1.SecretVolumeSource{
						// secret name must have "-token" suffix to be ignored (see lp:1925721)
						SecretName: "charm-data-token",
					},
				},
			},
			corev1.Volume{
				Name: "vol-projected",
				VolumeSource: corev1.VolumeSource{
					Projected: &corev1.ProjectedVolumeSource{},
				},
			},
			corev1.Volume{
				Name: "vol-configmap",
				VolumeSource: corev1.VolumeSource{
					ConfigMap: &corev1.ConfigMapVolumeSource{},
				},
			},
			corev1.Volume{
				Name: "vol-hostpath",
				VolumeSource: corev1.VolumeSource{
					HostPath: &corev1.HostPathVolumeSource{},
				},
			},
			corev1.Volume{
				Name: "vol-emptydir",
				VolumeSource: corev1.VolumeSource{
					EmptyDir: &corev1.EmptyDirVolumeSource{},
				},
			},
		)
		podSpec.Containers[0].VolumeMounts = append(podSpec.Containers[0].VolumeMounts,
			corev1.VolumeMount{Name: "vol-secret", MountPath: "path/secret"},
			corev1.VolumeMount{Name: "vol-projected", MountPath: "path/projected"},
			corev1.VolumeMount{Name: "vol-configmap", MountPath: "path/configmap"},
			corev1.VolumeMount{Name: "vol-hostpath", MountPath: "path/hostpath"},
			corev1.VolumeMount{Name: "vol-emptydir", MountPath: "path/emptydir"},
		)
		pod := corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Namespace:   s.namespace,
				Name:        fmt.Sprintf("%s-%d", s.appName, i),
				Labels:      map[string]string{"app.kubernetes.io/name": "gitlab"},
				Annotations: map[string]string{"juju.is/version": "2.9.37"},
			},
			Spec: podSpec,
			Status: corev1.PodStatus{
				PodIP: fmt.Sprintf("10.10.10.%d", i),
			},
		}
		switch i {
		case 0:
			pod.Status.Conditions = []corev1.PodCondition{
				{
					Type:    corev1.PodScheduled,
					Status:  corev1.ConditionFalse,
					Reason:  corev1.PodReasonUnschedulable,
					Message: "not enough resources",
				},
			}
		case 1:
			pod.Status.Conditions = []corev1.PodCondition{
				{
					Type:    corev1.PodScheduled,
					Status:  corev1.ConditionFalse,
					Reason:  "waiting",
					Message: "waiting to be scheduled",
				},
			}
		case 2:
			pod.DeletionTimestamp = &metav1.Time{
				Time: time.Now(),
			}
		case 3:
			pod.Status.Conditions = []corev1.PodCondition{}

		case 4:
			pod.Status.Conditions = []corev1.PodCondition{
				{
					Type:   corev1.PodScheduled,
					Status: corev1.ConditionTrue,
				},
				{
					Type:    corev1.PodInitialized,
					Status:  corev1.ConditionFalse,
					Reason:  resources.PodReasonContainersNotInitialized,
					Message: "initializing containers",
				},
			}
		case 5:
			pod.Status.Conditions = []corev1.PodCondition{
				{
					Type:   corev1.PodScheduled,
					Status: corev1.ConditionTrue,
				},
				{
					Type:    corev1.PodInitialized,
					Status:  corev1.ConditionFalse,
					Reason:  resources.PodReasonInitializing,
					Message: "initializing containers",
				},
			}
		case 6:
			pod.Status.Conditions = []corev1.PodCondition{
				{
					Type:   corev1.PodScheduled,
					Status: corev1.ConditionTrue,
				},
				{
					Type:    corev1.PodInitialized,
					Status:  corev1.ConditionFalse,
					Reason:  resources.PodReasonContainersNotInitialized,
					Message: "initializing containers",
				},
			}
			pod.Status.InitContainerStatuses = []corev1.ContainerStatus{
				{
					Name: "test-init-container",
					State: corev1.ContainerState{
						Waiting: &corev1.ContainerStateWaiting{
							Reason:  resources.PodReasonCrashLoopBackoff,
							Message: "I am broken",
						},
					},
				},
			}
		case 7:
			pod.Status.Conditions = []corev1.PodCondition{
				{
					Type:   corev1.PodScheduled,
					Status: corev1.ConditionTrue,
				},
				{
					Type:    corev1.ContainersReady,
					Status:  corev1.ConditionFalse,
					Reason:  resources.PodReasonContainersNotReady,
					Message: "starting containers",
				},
			}
			pod.Status.ContainerStatuses = []corev1.ContainerStatus{
				{
					Name: "test-container",
					State: corev1.ContainerState{
						Waiting: &corev1.ContainerStateWaiting{
							Reason:  "bad-reason",
							Message: "I am broken",
						},
					},
				},
			}
		case 8:
			pod.Status.Conditions = []corev1.PodCondition{
				{
					Type:   corev1.PodScheduled,
					Status: corev1.ConditionTrue,
				},
				{
					Type:   corev1.ContainersReady,
					Status: corev1.ConditionTrue,
				},
				{
					Type:   corev1.PodReady,
					Status: corev1.ConditionTrue,
				},
			}
		}
		_, err := s.client.CoreV1().Pods(s.namespace).Create(c.Context(), &pod, metav1.CreateOptions{})
		c.Assert(err, tc.ErrorIsNil)

		pvc := corev1.PersistentVolumeClaim{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: s.namespace,
				Name:      fmt.Sprintf("gitlab-database-appuuid-gitlab-%d", i),
			},
			Spec: corev1.PersistentVolumeClaimSpec{
				AccessModes: []corev1.PersistentVolumeAccessMode{
					corev1.ReadWriteOnce,
				},
				Resources: corev1.VolumeResourceRequirements{
					Requests: corev1.ResourceList{
						"storage": k8sresource.MustParse("1Gi"),
					},
				},
				VolumeName: fmt.Sprintf("pv-%d", i),
			},
			Status: corev1.PersistentVolumeClaimStatus{
				AccessModes: []corev1.PersistentVolumeAccessMode{
					corev1.ReadWriteOnce,
				},
				Capacity: corev1.ResourceList{
					"storage": k8sresource.MustParse("1Gi"),
				},
				Phase: corev1.ClaimBound,
			},
		}
		_, err = s.client.CoreV1().PersistentVolumeClaims(s.namespace).Create(c.Context(), &pvc, metav1.CreateOptions{})
		c.Assert(err, tc.ErrorIsNil)

		pv := corev1.PersistentVolume{
			ObjectMeta: metav1.ObjectMeta{
				Name: fmt.Sprintf("pv-%d", i),
			},
			Spec: corev1.PersistentVolumeSpec{
				AccessModes: []corev1.PersistentVolumeAccessMode{
					corev1.ReadWriteOnce,
				},
				Capacity: corev1.ResourceList{
					"storage": k8sresource.MustParse("1Gi"),
				},
			},
			Status: corev1.PersistentVolumeStatus{
				Phase:   corev1.VolumeBound,
				Message: "volume bound",
			},
		}
		_, err = s.client.CoreV1().PersistentVolumes().Create(c.Context(), &pv, metav1.CreateOptions{})
		c.Assert(err, tc.ErrorIsNil)
	}

	units, err := app.Units()
	c.Assert(err, tc.ErrorIsNil)

	mc := tc.NewMultiChecker()
	mc.AddExpr(`_[_].Status.Since`, tc.Ignore)
	mc.AddExpr(`_[_].FilesystemInfo[_].Status.Since`, tc.Ignore)
	mc.AddExpr(`_[_].FilesystemInfo[_].Volume.Status.Since`, tc.Ignore)

	c.Assert(units, mc, []caas.Unit{
		{
			Id:       "gitlab-0",
			Address:  "10.10.10.0",
			Ports:    []string(nil),
			Dying:    false,
			Stateful: true,
			Status: status.StatusInfo{
				Status:  "blocked",
				Message: "not enough resources",
			},
			FilesystemInfo: []caas.FilesystemInfo{
				{
					StorageName:  "gitlab-database",
					FilesystemId: "",
					Size:         1024,
					MountPoint:   "path/to/here",
					ReadOnly:     false,
					Status: status.StatusInfo{
						Status: "attached",
					},
					Volume: caas.VolumeInfo{
						VolumeId:   "pv-0",
						Size:       1024,
						Persistent: true,
						Status: status.StatusInfo{
							Status:  "attached",
							Message: "volume bound",
						},
					},
				},
			},
		},
		{
			Id:       "gitlab-1",
			Address:  "10.10.10.1",
			Ports:    []string(nil),
			Dying:    false,
			Stateful: true,
			Status: status.StatusInfo{
				Status:  "allocating",
				Message: "waiting to be scheduled",
			},
			FilesystemInfo: []caas.FilesystemInfo{
				{
					StorageName:  "gitlab-database",
					FilesystemId: "",
					Size:         1024,
					MountPoint:   "path/to/here",
					ReadOnly:     false,
					Status: status.StatusInfo{
						Status: "attached",
					},
					Volume: caas.VolumeInfo{
						VolumeId:   "pv-1",
						Size:       1024,
						Persistent: true,
						Status: status.StatusInfo{
							Status:  "attached",
							Message: "volume bound",
						},
					},
				},
			},
		},
		{
			Id:       "gitlab-2",
			Address:  "10.10.10.2",
			Ports:    []string(nil),
			Dying:    true,
			Stateful: true,
			Status: status.StatusInfo{
				Status: "terminated",
			},
			FilesystemInfo: []caas.FilesystemInfo{
				{
					StorageName:  "gitlab-database",
					FilesystemId: "",
					Size:         1024,
					MountPoint:   "path/to/here",
					ReadOnly:     false,
					Status: status.StatusInfo{
						Status: "attached",
					},
					Volume: caas.VolumeInfo{
						VolumeId:   "pv-2",
						Size:       1024,
						Persistent: true,
						Status: status.StatusInfo{
							Status:  "attached",
							Message: "volume bound",
						},
					},
				},
			},
		},
		{
			Id:       "gitlab-3",
			Address:  "10.10.10.3",
			Ports:    []string(nil),
			Dying:    false,
			Stateful: true,
			Status: status.StatusInfo{
				Status: "unknown",
			},
			FilesystemInfo: []caas.FilesystemInfo{
				{
					StorageName:  "gitlab-database",
					FilesystemId: "",
					Size:         1024,
					MountPoint:   "path/to/here",
					ReadOnly:     false,
					Status: status.StatusInfo{
						Status: "attached",
					},
					Volume: caas.VolumeInfo{
						VolumeId:   "pv-3",
						Size:       1024,
						Persistent: true,
						Status: status.StatusInfo{
							Status:  "attached",
							Message: "volume bound",
						},
					},
				},
			},
		},
		{
			Id:       "gitlab-4",
			Address:  "10.10.10.4",
			Ports:    []string(nil),
			Dying:    false,
			Stateful: true,
			Status: status.StatusInfo{
				Status:  "maintenance",
				Message: "initializing containers",
			},
			FilesystemInfo: []caas.FilesystemInfo{
				{
					StorageName:  "gitlab-database",
					FilesystemId: "",
					Size:         1024,
					MountPoint:   "path/to/here",
					ReadOnly:     false,
					Status: status.StatusInfo{
						Status: "attached",
					},
					Volume: caas.VolumeInfo{
						VolumeId:   "pv-4",
						Size:       1024,
						Persistent: true,
						Status: status.StatusInfo{
							Status:  "attached",
							Message: "volume bound",
						},
					},
				},
			},
		},
		{
			Id:       "gitlab-5",
			Address:  "10.10.10.5",
			Ports:    []string(nil),
			Dying:    false,
			Stateful: true,
			Status: status.StatusInfo{
				Status:  "maintenance",
				Message: "initializing containers",
			},
			FilesystemInfo: []caas.FilesystemInfo{
				{
					StorageName:  "gitlab-database",
					FilesystemId: "",
					Size:         1024,
					MountPoint:   "path/to/here",
					ReadOnly:     false,
					Status: status.StatusInfo{
						Status: "attached",
					},
					Volume: caas.VolumeInfo{
						VolumeId:   "pv-5",
						Size:       1024,
						Persistent: true,
						Status: status.StatusInfo{
							Status:  "attached",
							Message: "volume bound",
						},
					},
				},
			},
		},
		{
			Id:       "gitlab-6",
			Address:  "10.10.10.6",
			Ports:    []string(nil),
			Dying:    false,
			Stateful: true,
			Status: status.StatusInfo{
				Status:  "error",
				Message: "crash loop backoff: I am broken",
			},
			FilesystemInfo: []caas.FilesystemInfo{
				{
					StorageName:  "gitlab-database",
					FilesystemId: "",
					Size:         1024,
					MountPoint:   "path/to/here",
					ReadOnly:     false,
					Status: status.StatusInfo{
						Status: "attached",
					},
					Volume: caas.VolumeInfo{
						VolumeId:   "pv-6",
						Size:       1024,
						Persistent: true,
						Status: status.StatusInfo{
							Status:  "attached",
							Message: "volume bound",
						},
					},
				},
			},
		},
		{
			Id:       "gitlab-7",
			Address:  "10.10.10.7",
			Ports:    []string(nil),
			Dying:    false,
			Stateful: true,
			Status: status.StatusInfo{
				Status:  "error",
				Message: "unknown container reason \"bad-reason\": I am broken",
			},
			FilesystemInfo: []caas.FilesystemInfo{
				{
					StorageName:  "gitlab-database",
					FilesystemId: "",
					Size:         1024,
					MountPoint:   "path/to/here",
					ReadOnly:     false,
					Status: status.StatusInfo{
						Status: "attached",
					},
					Volume: caas.VolumeInfo{
						VolumeId:   "pv-7",
						Size:       1024,
						Persistent: true,
						Status: status.StatusInfo{
							Status:  "attached",
							Message: "volume bound",
						},
					},
				},
			},
		},
		{
			Id:       "gitlab-8",
			Address:  "10.10.10.8",
			Ports:    []string(nil),
			Dying:    false,
			Stateful: true,
			Status: status.StatusInfo{
				Status: "running",
			},
			FilesystemInfo: []caas.FilesystemInfo{
				{
					StorageName:  "gitlab-database",
					FilesystemId: "",
					Size:         1024,
					MountPoint:   "path/to/here",
					ReadOnly:     false,
					Status: status.StatusInfo{
						Status: "attached",
					},
					Volume: caas.VolumeInfo{
						VolumeId:   "pv-8",
						Size:       1024,
						Persistent: true,
						Status: status.StatusInfo{
							Status:  "attached",
							Message: "volume bound",
						},
					},
				},
			},
		},
	})
}

func (s *applicationSuite) TestServiceActive(c *tc.C) {
	app, _ := s.getApp(c, caas.DeploymentStateful, false)
	s.assertEnsure(
		c, app, false, constraints.Value{}, false, false, "", func() {},
	)
	defer s.assertDelete(c, app)

	testSvc := getDefaultSvc()
	testSvc.UID = "deadbeaf"
	testSvc.Spec.ClusterIP = "10.6.6.6"
	_, err := s.client.CoreV1().Services("test").Update(c.Context(), testSvc, metav1.UpdateOptions{})
	c.Assert(err, tc.ErrorIsNil)

	ss, err := s.client.AppsV1().StatefulSets("test").Get(c.Context(), "gitlab", metav1.GetOptions{})
	c.Assert(err, tc.ErrorIsNil)
	ss.Status.ReadyReplicas = 3
	_, err = s.client.AppsV1().StatefulSets("test").Update(c.Context(), ss, metav1.UpdateOptions{})
	c.Assert(err, tc.ErrorIsNil)

	svc, err := app.Service()
	c.Assert(err, tc.ErrorIsNil)

	since := time.Time{}
	c.Assert(svc, tc.DeepEquals, &caas.Service{
		Id: "deadbeaf",
		Addresses: network.ProviderAddresses{{
			MachineAddress: network.MachineAddress{
				Value: "10.6.6.6",
				Type:  "ipv4",
				Scope: "local-cloud",
			},
		}},
		Status: status.StatusInfo{
			Status: "active",
			Since:  &since,
		},
	})
}

func (s *applicationSuite) TestServiceNotSupportedDaemon(c *tc.C) {
	app, _ := s.getApp(c, caas.DeploymentDaemon, false)
	s.assertEnsure(
		c, app, false, constraints.Value{}, false, false, "", func() {},
	)
	defer s.assertDelete(c, app)

	testSvc := getDefaultSvc()
	testSvc.UID = "deadbeaf"
	testSvc.Spec.ClusterIP = "10.6.6.6"
	_, err := s.client.CoreV1().Services("test").Update(c.Context(), testSvc, metav1.UpdateOptions{})
	c.Assert(err, tc.ErrorIsNil)

	_, err = app.Service()
	c.Assert(err, tc.ErrorMatches, `deployment type "daemon" not supported`)
}

func (s *applicationSuite) TestServiceNotSupportedStateless(c *tc.C) {
	app, _ := s.getApp(c, caas.DeploymentStateless, false)
	s.assertEnsure(
		c, app, false, constraints.Value{}, false, false, "", func() {},
	)
	defer s.assertDelete(c, app)

	testSvc := getDefaultSvc()
	testSvc.UID = "deadbeaf"
	testSvc.Spec.ClusterIP = "10.6.6.6"
	_, err := s.client.CoreV1().Services("test").Update(c.Context(), testSvc, metav1.UpdateOptions{})
	c.Assert(err, tc.ErrorIsNil)

	_, err = app.Service()
	c.Assert(err, tc.ErrorMatches, `deployment type "stateless" not supported`)
}

func (s *applicationSuite) TestServiceTerminated(c *tc.C) {
	app, _ := s.getApp(c, caas.DeploymentStateful, false)
	s.assertEnsure(
		c, app, false, constraints.Value{}, false, false, "", func() {},
	)
	defer s.assertDelete(c, app)

	testSvc := getDefaultSvc()
	testSvc.UID = "deadbeaf"
	testSvc.Spec.ClusterIP = "10.6.6.6"
	_, err := s.client.CoreV1().Services("test").Update(c.Context(), testSvc, metav1.UpdateOptions{})
	c.Assert(err, tc.ErrorIsNil)

	ss, err := s.client.AppsV1().StatefulSets("test").Get(c.Context(), "gitlab", metav1.GetOptions{})
	c.Assert(err, tc.ErrorIsNil)
	now := metav1.Now()
	ss.DeletionTimestamp = &now
	_, err = s.client.AppsV1().StatefulSets("test").Update(c.Context(), ss, metav1.UpdateOptions{})
	c.Assert(err, tc.ErrorIsNil)

	svc, err := app.Service()
	c.Assert(err, tc.ErrorIsNil)

	since := time.Time{}
	c.Assert(svc, tc.DeepEquals, &caas.Service{
		Id: "deadbeaf",
		Addresses: network.ProviderAddresses{{
			MachineAddress: network.MachineAddress{
				Value: "10.6.6.6",
				Type:  "ipv4",
				Scope: "local-cloud",
			},
		}},
		Status: status.StatusInfo{
			Status: "terminated",
			Since:  &since,
		},
	})
}

func (s *applicationSuite) TestServiceError(c *tc.C) {
	app, _ := s.getApp(c, caas.DeploymentStateful, false)
	s.assertEnsure(
		c, app, false, constraints.Value{}, false, false, "", func() {},
	)
	defer s.assertDelete(c, app)

	testSvc := getDefaultSvc()
	testSvc.UID = "deadbeaf"
	testSvc.Spec.ClusterIP = "10.6.6.6"
	_, err := s.client.CoreV1().Services("test").Update(c.Context(), testSvc, metav1.UpdateOptions{})
	c.Assert(err, tc.ErrorIsNil)

	ss, err := s.client.AppsV1().StatefulSets("test").Get(c.Context(), "gitlab", metav1.GetOptions{})
	c.Assert(err, tc.ErrorIsNil)
	ss.Status.ReadyReplicas = 0
	_, err = s.client.AppsV1().StatefulSets("test").Update(c.Context(), ss, metav1.UpdateOptions{})
	c.Assert(err, tc.ErrorIsNil)

	evt := corev1.Event{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "test",
			Name:      "evt1",
		},
		InvolvedObject: corev1.ObjectReference{
			Name: "gitlab",
			Kind: "StatefulSet",
		},
		Type:    corev1.EventTypeWarning,
		Reason:  "FailedCreate",
		Message: "0/1 nodes are available: 1 pod has unbound immediate PersistentVolumeClaims.",
	}
	_, err = s.client.CoreV1().Events("test").Create(c.Context(), &evt, metav1.CreateOptions{})
	c.Assert(err, tc.ErrorIsNil)
	defer func() {
		_ = s.client.CoreV1().Events("test").Delete(c.Context(), evt.GetName(), metav1.DeleteOptions{})
	}()

	svc, err := app.Service()
	c.Assert(err, tc.ErrorIsNil)

	since := time.Time{}
	c.Assert(svc, tc.DeepEquals, &caas.Service{
		Id: "deadbeaf",
		Addresses: network.ProviderAddresses{{
			MachineAddress: network.MachineAddress{
				Value: "10.6.6.6",
				Type:  "ipv4",
				Scope: "local-cloud",
			},
		}},
		Status: status.StatusInfo{
			Status:  "error",
			Since:   &since,
			Message: "0/1 nodes are available: 1 pod has unbound immediate PersistentVolumeClaims.",
		},
	})
}

func (s *applicationSuite) TestEnsureConstraints(c *tc.C) {
	app, _ := s.getApp(c, caas.DeploymentStateful, false)
	s.assertEnsure(
		c, app, false, constraints.MustParse("mem=1G cpu-power=1000 arch=arm64"), true, false, "", func() {
			svc, err := s.client.CoreV1().Services("test").Get(c.Context(), "gitlab-endpoints", metav1.GetOptions{})
			c.Assert(err, tc.ErrorIsNil)
			c.Assert(svc, tc.DeepEquals, &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "gitlab-endpoints",
					Namespace: "test",
					Labels: map[string]string{
						"app.kubernetes.io/name":       "gitlab",
						"app.kubernetes.io/managed-by": "juju",
					},
					Annotations: map[string]string{
						"juju.is/version": "3.5-beta1",
						"service.alpha.kubernetes.io/tolerate-unready-endpoints": "true",
					},
				},
				Spec: corev1.ServiceSpec{
					Selector:                 map[string]string{"app.kubernetes.io/name": "gitlab"},
					Type:                     corev1.ServiceTypeClusterIP,
					ClusterIP:                "None",
					PublishNotReadyAddresses: true,
				},
			})

			ps := getPodSpec31()
			ps.NodeSelector = map[string]string{
				"kubernetes.io/arch": "arm64",
			}
			resourceRequests := corev1.ResourceList{
				corev1.ResourceCPU:    k8sresource.MustParse("1000m"),
				corev1.ResourceMemory: k8sresource.MustParse("1024Mi"),
			}
			ps.Containers[0].Resources.Requests = resourceRequests
			for i := range ps.Containers {
				ps.Containers[i].Resources.Limits = resourceRequests
			}

			ss, err := s.client.AppsV1().StatefulSets("test").Get(c.Context(), "gitlab", metav1.GetOptions{})
			c.Assert(err, tc.ErrorIsNil)
			c.Assert(ss, tc.DeepEquals, &appsv1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "gitlab",
					Namespace: "test",
					Labels: map[string]string{
						"app.kubernetes.io/name":       "gitlab",
						"app.kubernetes.io/managed-by": "juju",
					},
					Annotations: map[string]string{
						"juju.is/version":  "3.5-beta1",
						"app.juju.is/uuid": "appuuid",
					},
				},
				Spec: appsv1.StatefulSetSpec{
					Replicas: pointer.Int32Ptr(3),
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"app.kubernetes.io/name": "gitlab",
						},
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels:      map[string]string{"app.kubernetes.io/name": "gitlab"},
							Annotations: map[string]string{"juju.is/version": "3.5-beta1"},
						},
						Spec: ps,
					},
					VolumeClaimTemplates: []corev1.PersistentVolumeClaim{
						{
							ObjectMeta: metav1.ObjectMeta{
								Name: "gitlab-database-appuuid",
								Labels: map[string]string{
									"storage.juju.is/name":         "database",
									"app.kubernetes.io/managed-by": "juju",
								},
								Annotations: map[string]string{
									"foo":                  "bar",
									"storage.juju.is/name": "database",
								}},
							Spec: corev1.PersistentVolumeClaimSpec{
								StorageClassName: pointer.StringPtr("test-workload-storage"),
								AccessModes:      []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
								Resources: corev1.VolumeResourceRequirements{
									Requests: corev1.ResourceList{
										corev1.ResourceStorage: k8sresource.MustParse("100Mi"),
									},
								},
							},
						},
					},
					PodManagementPolicy: appsv1.ParallelPodManagement,
					ServiceName:         "gitlab-endpoints",
				},
			})
		},
	)
}

func (s *applicationSuite) TestPullSecretUpdate(c *tc.C) {
	app, _ := s.getApp(c, caas.DeploymentStateful, false)

	unusedPullSecret := corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "gitlab-oldcontainer-secret",
			Namespace: "test",
			Labels: map[string]string{
				"app.kubernetes.io/name":       "gitlab",
				"app.kubernetes.io/managed-by": "juju",
			},
			Annotations: map[string]string{"juju.is/version": "3.5-beta1"},
		},
		Type: corev1.SecretTypeDockerConfigJson,
		Data: map[string][]byte{
			corev1.DockerConfigJsonKey: []byte("wow"),
		},
	}

	_, err := s.client.CoreV1().Secrets(s.namespace).Create(c.Context(), &unusedPullSecret,
		metav1.CreateOptions{})
	c.Assert(err, tc.ErrorIsNil)

	pullSecretConfig, _ := k8sutils.CreateDockerConfigJSON("username-old", "password-old", "docker.io/library/nginx:latest")
	nginxPullSecret := corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "gitlab-nginx-secret",
			Namespace: "test",
			Labels: map[string]string{
				"app.kubernetes.io/name":       "gitlab",
				"app.kubernetes.io/managed-by": "juju",
			},
			Annotations: map[string]string{"juju.is/version": "3.5-beta1"},
		},
		Type: corev1.SecretTypeDockerConfigJson,
		Data: map[string][]byte{
			corev1.DockerConfigJsonKey: pullSecretConfig,
		},
	}
	_, err = s.client.CoreV1().Secrets(s.namespace).Create(c.Context(), &nginxPullSecret,
		metav1.CreateOptions{})
	c.Assert(err, tc.ErrorIsNil)

	s.assertEnsure(c, app, false, constraints.Value{}, true, false, "", func() {})

	_, err = s.client.CoreV1().Secrets(s.namespace).Get(c.Context(), "gitlab-oldcontainer-secret", metav1.GetOptions{})
	c.Assert(err, tc.ErrorMatches, `secrets "gitlab-oldcontainer-secret" not found`)

	secret, err := s.client.CoreV1().Secrets(s.namespace).Get(c.Context(), "gitlab-nginx-secret", metav1.GetOptions{})
	c.Assert(err, tc.ErrorIsNil)
	c.Assert(secret, tc.NotNil)
	newPullSecretConfig, _ := k8sutils.CreateDockerConfigJSON("username", "password", "docker.io/library/nginx:latest")
	newNginxPullSecret := nginxPullSecret
	newNginxPullSecret.Data = map[string][]byte{
		corev1.DockerConfigJsonKey: newPullSecretConfig,
	}
	c.Assert(*secret, tc.DeepEquals, newNginxPullSecret)
}

func (s *applicationSuite) TestPVCNames(c *tc.C) {
	claims := []*corev1.PersistentVolumeClaim{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "storage_a-abcd1234-gitlab-0",
				Namespace: "test",
				Labels: map[string]string{
					"app.kubernetes.io/managed-by": "juju",
					"app.kubernetes.io/name":       "gitlab",
					"storage.juju.is/name":         "storage_a",
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "gitlab-storage_b-abcd1234-gitlab-0",
				Namespace: "test",
				Labels: map[string]string{
					"app.kubernetes.io/managed-by": "juju",
					"app.kubernetes.io/name":       "gitlab",
					"storage.juju.is/name":         "storage_b",
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "gitlab-storage_g-abcd666-gitlab-0",
				Namespace: "test",
				Labels: map[string]string{
					"app.kubernetes.io/managed-by": "juju",
					"app.kubernetes.io/name":       "gitlab",
					"storage.juju.is/name":         "storage_g",
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "juju-storage_c-42",
				Namespace: "test",
				Labels: map[string]string{
					"app.kubernetes.io/managed-by": "juju",
					"app.kubernetes.io/name":       "gitlab",
					"storage.juju.is/name":         "storage_c",
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "storage_d-abcd1234-gitlab-0",
				Namespace: "test",
				Labels: map[string]string{
					"app.kubernetes.io/managed-by": "juju",
					"app.kubernetes.io/name":       "another-app",
					"storage.juju.is/name":         "storage_d",
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "storage_e-abcd1236-gitlab-0",
				Namespace: "test",
				Labels: map[string]string{
					"app.kubernetes.io/managed-by": "juju",
					"app.kubernetes.io/name":       "gitlab",
					// no "storage.juju.is/name" label -- will be ignored
				},
			},
		},
	}
	for _, claim := range claims {
		_, err := s.client.CoreV1().PersistentVolumeClaims("test").Create(c.Context(), claim, metav1.CreateOptions{})
		c.Assert(err, tc.ErrorIsNil)
	}

	names, err := application.PVCNames(s.client, "test", "gitlab", "abcd1234")
	c.Assert(err, tc.ErrorIsNil)
	c.Assert(names, tc.DeepEquals, map[string]string{
		"gitlab-storage_a": "storage_a-abcd1234",
		"gitlab-storage_b": "gitlab-storage_b-abcd1234",
		"gitlab-storage_c": "juju-storage_c-42",
	})
}

func (s *applicationSuite) TestLimits(c *tc.C) {
	limits := corev1.ResourceList{
		corev1.ResourceCPU:    *k8sresource.NewMilliQuantity(1000, k8sresource.DecimalSI),
		corev1.ResourceMemory: *k8sresource.NewQuantity(1024*1024*1024, k8sresource.BinarySI),
	}

	app, _ := s.getApp(c, caas.DeploymentStateful, false)
	s.assertEnsure(
		c, app, false, constraints.MustParse("mem=1G cpu-power=1000 arch=arm64"), true, false, "", func() {
			ss, err := s.client.AppsV1().StatefulSets("test").Get(c.Context(), "gitlab", metav1.GetOptions{})
			c.Assert(err, tc.ErrorIsNil)
			for _, ctr := range ss.Spec.Template.Spec.Containers {
				c.Check(ctr.Resources.Limits, tc.DeepEquals, limits)
			}
		},
	)
}

func (s *applicationSuite) TestEnsureUpdatedConstraints(c *tc.C) {
	app, _ := s.getApp(c, caas.DeploymentStateful, false)
	s.assertEnsure(
		c, app, false, constraints.MustParse("mem=1G cpu-power=1000"), true, true, "3.6.8", func() {
			ps := getPodSpec368()
			charmResourceMemRequest := corev1.ResourceList{
				corev1.ResourceMemory: k8sresource.MustParse(constants.CharmMemRequestMi),
			}
			charmResourceMemLimit := corev1.ResourceList{
				corev1.ResourceMemory: k8sresource.MustParse(constants.CharmMemLimitMi),
			}

			workloadResourceLimits := corev1.ResourceList{
				corev1.ResourceCPU:    k8sresource.MustParse("1000m"),
				corev1.ResourceMemory: k8sresource.MustParse("1024Mi"),
			}

			for i, container := range ps.Containers {
				if container.Name == constants.ApplicationCharmContainer {
					continue
				}
				ps.Containers[i].Resources.Requests = workloadResourceLimits
				ps.Containers[i].Resources.Limits = workloadResourceLimits
			}
			ss, err := s.client.AppsV1().StatefulSets("test").Get(context.TODO(), "gitlab", metav1.GetOptions{})
			c.Assert(err, tc.ErrorIsNil)
			for _, ctr := range ss.Spec.Template.Spec.Containers {
				if ctr.Name == constants.ApplicationCharmContainer {
					c.Assert(ctr.Resources.Requests, tc.DeepEquals, charmResourceMemRequest)
					c.Assert(ctr.Resources.Limits, tc.DeepEquals, charmResourceMemLimit)
					continue
				}

				c.Check(ctr.Resources.Requests.Cpu().Equal(*workloadResourceLimits.Cpu()), tc.IsTrue)
				c.Check(ctr.Resources.Requests.Memory().Equal(*workloadResourceLimits.Memory()), tc.IsTrue)

				c.Check(ctr.Resources.Requests.Cpu().Equal(*workloadResourceLimits.Cpu()), tc.IsTrue)
				c.Check(ctr.Resources.Requests.Memory().Equal(*workloadResourceLimits.Memory()), tc.IsTrue)
			}
		},
	)
}

func int64Ptr(a int64) *int64 {
	return &a
}
