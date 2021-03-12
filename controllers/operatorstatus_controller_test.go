package controllers

import (
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab-ce.alauda.cn/micro-service/operator-monitor/api/v1alpha1"
	"gitlab-ce.alauda.cn/micro-service/operator-monitor/pkg/mock"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

var _ = Describe("Status Controller", func() {
	Context("get status ", func() {
		var got v1alpha1.OperatorStatus
		BeforeEach(func() {
			objs := mock.GetObjects()
			for _, obj := range objs {
				Expect(k8sClient.Create(context.Background(), obj)).Should(Succeed())
			}
			By("create asm status objects")
		})
		It("asm status", func() {
			key := client.ObjectKey{Name: "asm"}
			Eventually(func() bool {
				err := k8sClient.Get(context.Background(), key, &got)
				return err == nil
			}, time.Second*45, 1).Should(BeTrue())
			Expect(got.Name).Should(Equal(key.Name))
		})
	})
})
