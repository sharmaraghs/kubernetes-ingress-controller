package store

import (
	"testing"

	"github.com/go-logr/logr"
	"github.com/ngrok/kubernetes-ingress-controller/internal/errors"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	netv1 "k8s.io/api/networking/v1"
)

const ngrokIngressClass = "ngrok"
const defaultControllerName = "k8s.ngrok.com/ingress-controller"

func TestStore(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Store package Test Suite")
}

var _ = Describe("Store", func() {

	var store Storer
	BeforeEach(func() {
		// create a fake logger to pass into the cachestore
		logger := logr.New(logr.Discard().GetSink())
		cacheStores := NewCacheStores(logger)
		store = New(cacheStores, defaultControllerName, logger)
	})

	var _ = Describe("GetIngressClassV1", func() {
		Context("when the ingress class exists", func() {
			BeforeEach(func() {
				ic := NewTestIngressClass(ngrokIngressClass, true, true)
				_ = store.Add(&ic)
			})
			It("returns the ingress class", func() {
				ic, err := store.GetIngressClassV1(ngrokIngressClass)
				Expect(err).ToNot(HaveOccurred())
				Expect(ic.Name).To(Equal(ngrokIngressClass))
			})
		})
		Context("when the ingress class does not exist", func() {
			It("returns an error", func() {
				ic, err := store.GetIngressClassV1("does-not-exist")
				Expect(err).To(HaveOccurred())
				Expect(ic).To(BeNil())
			})
		})
	})

	var _ = Describe("GetIngressV1", func() {
		Context("when the ingress exists", func() {
			BeforeEach(func() {
				ing := NewTestIngressV1("test-ingress", "test-namespace")
				_ = store.Add(&ing)
			})
			It("returns the ingress", func() {
				ing, err := store.GetIngressV1("test-ingress", "test-namespace")
				Expect(err).ToNot(HaveOccurred())
				Expect(ing.Name).To(Equal("test-ingress"))
			})
		})
		Context("when the ingress does not exist", func() {
			It("returns an error", func() {
				ing, err := store.GetIngressV1("does-not-exist", "does-not-exist")
				Expect(err).To(HaveOccurred())
				Expect(ing).To(BeNil())
			})
		})
	})

	var _ = Describe("GetServiceV1", func() {
		Context("when the service exists", func() {
			BeforeEach(func() {
				svc := NewTestServiceV1("test-service", "test-namespace")
				_ = store.Add(&svc)
			})
			It("returns the service", func() {
				svc, err := store.GetServiceV1("test-service", "test-namespace")
				Expect(err).ToNot(HaveOccurred())
				Expect(svc.Name).To(Equal("test-service"))
			})
		})
		Context("when the service does not exist", func() {
			It("returns an error", func() {
				svc, err := store.GetServiceV1("does-not-exist", "does-not-exist")
				Expect(err).To(HaveOccurred())
				Expect(svc).To(BeNil())
			})
		})
	})

	var _ = Describe("GetNgrokIngressV1", func() {
		Context("when the ngrok ingress exists", func() {
			BeforeEach(func() {
				ing := NewTestIngressV1WithClass("test-ingress", "test-namespace", ngrokIngressClass)
				_ = store.Add(&ing)
				ic := NewTestIngressClass(ngrokIngressClass, true, true)
				_ = store.Add(&ic)
			})
			It("returns the ngrok ingress", func() {
				ing, err := store.GetNgrokIngressV1("test-ingress", "test-namespace")
				Expect(err).ToNot(HaveOccurred())
				Expect(ing.Name).To(Equal("test-ingress"))
			})
			It("Filters out ingresses that don't match the ngrok ingress class", func() {
				ingNotNgrok := NewTestIngressV1WithClass("ingNotNgrok", "test-namespace", "not-ngrok")
				_ = store.Add(&ingNotNgrok)

				ing, err := store.GetNgrokIngressV1("ingNotNgrok", "test-namespace")
				Expect(err).To(HaveOccurred())
				Expect(ing).To(BeNil())
			})
			It("Filters finds ones without a class if we are default", func() {
				ingNoClass := NewTestIngressV1("ingNoClass", "test-namespace")
				_ = store.Add(&ingNoClass)

				ing, err := store.GetNgrokIngressV1("ingNoClass", "test-namespace")
				Expect(err).ToNot(HaveOccurred())
				Expect(ing.Name).To(Equal("ingNoClass"))
			})
		})
		Context("when the ngrok ingress does not exist", func() {
			It("returns an error", func() {
				ing, err := store.GetNgrokIngressV1("does-not-exist", "does-not-exist")
				Expect(err).To(HaveOccurred())
				Expect(ing).To(BeNil())
			})
		})
	})

	var _ = Describe("ListNgrokIngressClassesV1", func() {
		Context("when there are ngrok ingress classes", func() {
			BeforeEach(func() {
				ic1 := NewTestIngressClass("ngrok1", true, true)
				_ = store.Add(&ic1)
				ic2 := NewTestIngressClass("ngrok2", true, true)
				_ = store.Add(&ic2)
				ic3 := NewTestIngressClass("different", true, false)
				_ = store.Add(&ic3)
			})
			It("returns the ngrok ingress classes and doesn't return the different one", func() {
				ics := store.ListNgrokIngressClassesV1()
				Expect(len(ics)).To(Equal(2))
			})
		})
		Context("when there are no ngrok ingress classes", func() {
			It("doesn't error", func() {
				ics := store.ListNgrokIngressClassesV1()
				Expect(len(ics)).To(Equal(0))
			})
		})
	})

	var _ = Describe("ListNgrokIngressesV1", func() {
		icUsDefault := NewTestIngressClass("ngrok", true, true)
		icUsNotDefault := NewTestIngressClass("ngrok", false, true)
		icOtherDefault := NewTestIngressClass("test", true, false)
		icOtherNotDefault := NewTestIngressClass("test", false, false)

		var _ = DescribeTable("IngressClassFiltering", func(ingressClasses []netv1.IngressClass, expectedMatchingIngressesCount int) {
			iMatching := NewTestIngressV1WithClass("test1", "test", "ngrok")
			iNotMatching := NewTestIngressV1WithClass("test2", "test", "test")
			iNoClass := NewTestIngressV1("test3", "test")
			_ = store.Add(&iMatching)
			_ = store.Add(&iNotMatching)
			_ = store.Add(&iNoClass)
			for _, ic := range ingressClasses {
				_ = store.Add(&ic)
			}
			ings := store.ListNgrokIngressesV1()
			Expect(len(ings)).To(Equal(expectedMatchingIngressesCount))
		},
			Entry("No ingress classes", []netv1.IngressClass{}, 0),
			Entry("just us not as default", []netv1.IngressClass{icUsNotDefault}, 1),
			Entry("just us as default", []netv1.IngressClass{icUsDefault}, 2),
			Entry("just another not as default", []netv1.IngressClass{icOtherNotDefault}, 0),
			Entry("just another as default", []netv1.IngressClass{icOtherDefault}, 0),
			Entry("us and another neither default", []netv1.IngressClass{icUsNotDefault, icOtherNotDefault}, 1),
			Entry("us and another them default", []netv1.IngressClass{icUsNotDefault, icOtherDefault}, 1),
			Entry("us and another us default", []netv1.IngressClass{icUsDefault, icOtherNotDefault}, 2),
			Entry("us and another both default", []netv1.IngressClass{icUsDefault, icOtherDefault}, 2),
		)
	})

	var _ = Describe("ListNgrokModulesV1", func() {
		Context("when there are NgrokModuleSets", func() {
			BeforeEach(func() {
				m1 := NewTestNgrokModuleSet("ngrok", "test", true)
				_ = store.Add(&m1)
				m2 := NewTestNgrokModuleSet("ngrok", "test2", true)
				_ = store.Add(&m2)
				m3 := NewTestNgrokModuleSet("test", "test", true)
				_ = store.Add(&m3)
			})
			It("returns the NgrokModuleSet", func() {
				modules := store.ListNgrokModuleSetsV1()
				Expect(len(modules)).To(Equal(3))
			})
		})
		Context("when there are no NgrokModuleSets", func() {
			It("doesn't error", func() {
				modules := store.ListNgrokModuleSetsV1()
				Expect(len(modules)).To(Equal(0))
			})
		})
	})

	var _ = Describe("GetNgrokModuleSetV1", func() {
		Context("when the NgrokModuleSet exists", func() {
			BeforeEach(func() {
				m := NewTestNgrokModuleSet("ngrok", "test", true)
				_ = store.Add(&m)
			})
			It("returns the NgrokModuleSet", func() {
				modset, err := store.GetNgrokModuleSetV1("ngrok", "test")
				Expect(err).ToNot(HaveOccurred())
				Expect(modset.Modules.Compression.Enabled).To(Equal(true))
			})
		})
		Context("when the NgrokModuleSet does not exist", func() {
			It("returns an error", func() {
				modset, err := store.GetNgrokModuleSetV1("does-not-exist", "does-not-exist")
				Expect(err).To(HaveOccurred())
				Expect(errors.IsErrorNotFound(err)).To(Equal(true))
				Expect(modset).To(BeNil())
			})
		})
	})
})
