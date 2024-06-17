package status_test

import (
	"github.com/awslabs/operatorpkg/status"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Conditions", func() {
	It("should correctly toggle conditions", func() {
		testObject := TestObject{}
		// Conditions should be initialized
		conditions := testObject.StatusConditions()
		Expect(conditions.Get(ConditionTypeFoo).GetStatus()).To(Equal(metav1.ConditionUnknown))
		Expect(conditions.Get(ConditionTypeBar).GetStatus()).To(Equal(metav1.ConditionUnknown))
		Expect(conditions.Root().GetStatus()).To(Equal(metav1.ConditionUnknown))
		// Update the condition to true
		Expect(conditions.SetTrue(ConditionTypeFoo)).To(BeTrue())
		fooCondition := conditions.Get(ConditionTypeFoo)
		Expect(fooCondition.Type).To(Equal(ConditionTypeFoo))
		Expect(fooCondition.Status).To(Equal(metav1.ConditionTrue))
		Expect(fooCondition.Reason).To(Equal(ConditionTypeFoo)) // default to type
		Expect(fooCondition.Message).To(Equal(""))              // default to type
		Expect(fooCondition.LastTransitionTime.UnixNano()).To(BeNumerically(">", 0))
		Expect(conditions.Root().GetStatus()).To(Equal(metav1.ConditionUnknown))
		time.Sleep(1 * time.Nanosecond)
		// Update the other condition to false
		Expect(conditions.SetFalse(ConditionTypeBar, "reason", "message")).To(BeTrue())
		fooCondition2 := conditions.Get(ConditionTypeBar)
		Expect(fooCondition2.Type).To(Equal(ConditionTypeBar))
		Expect(fooCondition2.Status).To(Equal(metav1.ConditionFalse))
		Expect(fooCondition2.Reason).To(Equal("reason"))
		Expect(fooCondition2.Message).To(Equal("message"))
		Expect(fooCondition.LastTransitionTime.UnixNano()).To(BeNumerically(">", 0))
		Expect(conditions.Root().GetStatus()).To(Equal(metav1.ConditionFalse))
		time.Sleep(1 * time.Nanosecond)
		// transition the root condition to true
		Expect(conditions.SetTrueWithReason(ConditionTypeBar, "reason", "message")).To(BeTrue())
		updatedFooCondition := conditions.Get(ConditionTypeBar)
		Expect(updatedFooCondition.Type).To(Equal(ConditionTypeBar))
		Expect(updatedFooCondition.Status).To(Equal(metav1.ConditionTrue))
		Expect(updatedFooCondition.Reason).To(Equal("reason"))
		Expect(updatedFooCondition.Message).To(Equal("message"))
		Expect(updatedFooCondition.LastTransitionTime.UnixNano()).To(BeNumerically(">", fooCondition.LastTransitionTime.UnixNano()))
		Expect(conditions.Root().GetStatus()).To(Equal(metav1.ConditionTrue))
		time.Sleep(1 * time.Nanosecond)
		// Transition if the status is the same, but the Reason is different
		Expect(conditions.SetFalse(ConditionTypeBar, "another-reason", "another-message")).To(BeTrue())
		updatedBarCondition := conditions.Get(ConditionTypeBar)
		Expect(updatedBarCondition.Type).To(Equal(ConditionTypeBar))
		Expect(updatedBarCondition.Status).To(Equal(metav1.ConditionFalse))
		Expect(updatedBarCondition.Reason).To(Equal("another-reason"))
		Expect(updatedBarCondition.LastTransitionTime.UnixNano()).ToNot(BeNumerically("==", fooCondition2.LastTransitionTime.UnixNano()))
		// Dont transition if reason and message are the same
		Expect(conditions.SetTrue(ConditionTypeFoo)).To(BeFalse())
		Expect(conditions.SetFalse(ConditionTypeBar, "another-reason", "another-message")).To(BeFalse())
	})

	It("all true", func() {
		testObject := TestObject{}
		Expect(testObject.StatusConditions().IsTrue()).To(BeTrue())
		Expect(testObject.StatusConditions().IsTrue(ConditionTypeFoo)).To(BeFalse())
		Expect(testObject.StatusConditions().IsTrue(ConditionTypeBar)).To(BeFalse())
		Expect(testObject.StatusConditions().IsTrue(ConditionTypeBaz)).To(BeFalse())
		Expect(testObject.StatusConditions().IsTrue(ConditionTypeFoo, ConditionTypeBar)).To(BeFalse())
		Expect(testObject.StatusConditions().IsTrue(ConditionTypeFoo, ConditionTypeBaz)).To(BeFalse())
		Expect(testObject.StatusConditions().IsTrue(ConditionTypeFoo, ConditionTypeBar, ConditionTypeBaz)).To(BeFalse())

		testObject.StatusConditions().SetFalse(ConditionTypeBaz, "reason", "message")
		Expect(testObject.StatusConditions().IsTrue()).To(BeTrue())
		Expect(testObject.StatusConditions().IsTrue(ConditionTypeFoo)).To(BeFalse())
		Expect(testObject.StatusConditions().IsTrue(ConditionTypeBar)).To(BeFalse())
		Expect(testObject.StatusConditions().IsTrue(ConditionTypeBaz)).To(BeFalse())
		Expect(testObject.StatusConditions().IsTrue(ConditionTypeFoo, ConditionTypeBar)).To(BeFalse())
		Expect(testObject.StatusConditions().IsTrue(ConditionTypeFoo, ConditionTypeBaz)).To(BeFalse())
		Expect(testObject.StatusConditions().IsTrue(ConditionTypeFoo, ConditionTypeBar, ConditionTypeBaz)).To(BeFalse())

		testObject.StatusConditions().SetTrue(ConditionTypeFoo)
		Expect(testObject.StatusConditions().IsTrue()).To(BeTrue())
		Expect(testObject.StatusConditions().IsTrue(ConditionTypeFoo)).To(BeTrue())
		Expect(testObject.StatusConditions().IsTrue(ConditionTypeBar)).To(BeFalse())
		Expect(testObject.StatusConditions().IsTrue(ConditionTypeBaz)).To(BeFalse())
		Expect(testObject.StatusConditions().IsTrue(ConditionTypeFoo, ConditionTypeBar)).To(BeFalse())
		Expect(testObject.StatusConditions().IsTrue(ConditionTypeFoo, ConditionTypeBaz)).To(BeFalse())
		Expect(testObject.StatusConditions().IsTrue(ConditionTypeFoo, ConditionTypeBar, ConditionTypeBaz)).To(BeFalse())

		testObject.StatusConditions().SetTrue(ConditionTypeBar)
		Expect(testObject.StatusConditions().IsTrue()).To(BeTrue())
		Expect(testObject.StatusConditions().IsTrue(ConditionTypeFoo)).To(BeTrue())
		Expect(testObject.StatusConditions().IsTrue(ConditionTypeBar)).To(BeTrue())
		Expect(testObject.StatusConditions().IsTrue(ConditionTypeBaz)).To(BeFalse())
		Expect(testObject.StatusConditions().IsTrue(ConditionTypeFoo, ConditionTypeBar)).To(BeTrue())
		Expect(testObject.StatusConditions().IsTrue(ConditionTypeFoo, ConditionTypeBaz)).To(BeFalse())
		Expect(testObject.StatusConditions().IsTrue(ConditionTypeFoo, ConditionTypeBar, ConditionTypeBaz)).To(BeFalse())

		testObject.StatusConditions().SetTrue(ConditionTypeBaz)
		Expect(testObject.StatusConditions().IsTrue()).To(BeTrue())
		Expect(testObject.StatusConditions().IsTrue(ConditionTypeFoo)).To(BeTrue())
		Expect(testObject.StatusConditions().IsTrue(ConditionTypeBar)).To(BeTrue())
		Expect(testObject.StatusConditions().IsTrue(ConditionTypeBaz)).To(BeTrue())
		Expect(testObject.StatusConditions().IsTrue(ConditionTypeFoo, ConditionTypeBar)).To(BeTrue())
		Expect(testObject.StatusConditions().IsTrue(ConditionTypeFoo, ConditionTypeBaz)).To(BeTrue())
		Expect(testObject.StatusConditions().IsTrue(ConditionTypeFoo, ConditionTypeBar, ConditionTypeBaz)).To(BeTrue())
	})
	It("should return nil when status conditions are not found on the expected path", func() {
		testObject := &unstructured.Unstructured{Object: map[string]interface{}{
			"invalid": map[string]interface{}{
				"invalid": "invalid",
			},
		}}
		obj, err := status.FromUnstructured(testObject)
		Expect(obj).To(BeNil())
		Expect(err).To(HaveOccurred())
	})
	It("should validate status condition on unstructured object status is false", func() {
		conditionObj, err := status.FromUnstructured(createUnstructuredStatusConditions("False"))
		Expect(err).To(BeNil())
		Expect(conditionObj).ToNot(BeNil())
		Expect(conditionObj.StatusConditions().IsTrue(status.ConditionReady)).To(BeFalse())
	})
	It("should validate status condition on unstructured object status is true", func() {
		conditionObj, err := status.FromUnstructured(createUnstructuredStatusConditions("True"))
		Expect(err).To(BeNil())
		Expect(conditionObj).ToNot(BeNil())
		Expect(conditionObj.StatusConditions().IsTrue(status.ConditionReady)).To(BeTrue())
	})
	It("should set condition on unstructured object", func() {
		testObject := &unstructured.Unstructured{Object: map[string]interface{}{
			"status": map[string]interface{}{
				"conditions": []interface{}{},
			},
		}}
		conditions := []status.Condition{
			{
				Type:    status.ConditionSucceeded,
				Status:  metav1.ConditionFalse,
				Reason:  "reason",
				Message: "message",
			},
		}
		conditionObj, err := status.FromUnstructured(testObject)
		Expect(err).ToNot(HaveOccurred())
		conditionObj.SetConditions(conditions)
		c, found, err := unstructured.NestedSlice(testObject.Object, "status", "conditions")
		Expect(err).To(BeNil())
		Expect(found).To(BeTrue())
		Expect(len(c)).To(BeEquivalentTo(1))
	})
})

func createUnstructuredStatusConditions(status string) *unstructured.Unstructured {
	return &unstructured.Unstructured{Object: map[string]interface{}{
		"status": map[string]interface{}{
			"conditions": []interface{}{
				map[string]interface{}{
					"type":    "Ready",
					"status":  status,
					"message": "message",
					"reason":  "reason",
				},
			},
		},
	}}
}
