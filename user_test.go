package user

import (
	"context"

	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/format"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestUserService(t *testing.T) {
	format.MaxLength = 0

	RegisterFailHandler(Fail)
	RunSpecs(t, "UserService")
}

var _ = Describe("UserService", func() {
	db := &gorm.DB{}
	db, _ = gorm.Open(postgres.New(postgres.Config{DSN: "host=192.168.10.138 port=5432 user=postgres password=postgres dbname=users_list sslmode=disable"}))
	var (
		app *ServiceUserImpl
		ctx context.Context
	)

	BeforeEach(func() {
		app = New(db)
		ctx = context.WithValue(context.Background(), "", app)
		ctx.Value("")

	})

	Context("Public functions", func() {
		When("add string Attr", func() {
			It("Success", func() {
				id, e := app.NewUser()
				Expect(e).ShouldNot(HaveOccurred())

				code := "name_t"

				e = app.UpdateUserStringAttr(id, code, "anton_t")
				Expect(e).ShouldNot(HaveOccurred())

			})
		})
		When("update string Attr", func() {
			It("Success", func() {
				id, e := app.NewUser()
				Expect(e).ShouldNot(HaveOccurred())

				code := "name_t"

				e = app.UpdateUserStringAttr(id, code, "maxim_t")
				Expect(e).ShouldNot(HaveOccurred())

				e = app.UpdateUserStringAttr(id, code, "max_t")
				Expect(e).ShouldNot(HaveOccurred())

				app.DeleteUserStringAttr(id, code)
			})
		})
		When("delete string Attr", func() {
			It("Success", func() {
				id, e := app.NewUser()
				Expect(e).ShouldNot(HaveOccurred())

				code := "name_t"

				e = app.UpdateUserStringAttr(id, code, "anderey_t")
				Expect(e).ShouldNot(HaveOccurred())

				e = app.DeleteUserStringAttr(id, code)
				Expect(e).ShouldNot(HaveOccurred())

				e = app.DeleteUserStringAttr(id, "not_existing_code")
				Expect(e).ShouldNot(HaveOccurred())

			})
		})
		When("get string Attr by user", func() {
			It("Success", func() {
				id1, e := app.NewUser()
				Expect(e).ShouldNot(HaveOccurred())
				id2, e := app.NewUser()
				Expect(e).ShouldNot(HaveOccurred())

				code1 := "name_t"
				code2 := "city_t"

				e = app.UpdateUserStringAttr(id1, code1, "anderey_t")
				Expect(e).ShouldNot(HaveOccurred())
				e = app.UpdateUserStringAttr(id1, code2, "moscow_t")
				Expect(e).ShouldNot(HaveOccurred())

				e = app.UpdateUserStringAttr(id2, code1, "anderey_t")
				Expect(e).ShouldNot(HaveOccurred())
				e = app.UpdateUserStringAttr(id2, code2, "piterburg_t")
				Expect(e).ShouldNot(HaveOccurred())
				e = app.DeleteUserStringAttr(id2, code2)
				Expect(e).ShouldNot(HaveOccurred())

				attrs, e := app.GetUserAttr(id1)
				Expect(e).ShouldNot(HaveOccurred())
				Expect(attrs).ShouldNot(BeNil())
				Expect(len(attrs)).Should(Equal(2))

				attrs, e = app.GetUserAttr(id2)
				Expect(e).ShouldNot(HaveOccurred())
				Expect(attrs).ShouldNot(BeNil())
				Expect(len(attrs)).Should(Equal(1))

				app.DeleteUserStringAttr(id1, code1)
				app.DeleteUserStringAttr(id1, code2)
				app.DeleteUserStringAttr(id2, code1)
				app.DeleteUserStringAttr(id2, code2)

			})
		})
		When("get users by string Attr", func() {
			It("Success", func() {
				id1, e := app.NewUser()
				Expect(e).ShouldNot(HaveOccurred())
				id2, e := app.NewUser()
				Expect(e).ShouldNot(HaveOccurred())
				id3, e := app.NewUser()
				Expect(e).ShouldNot(HaveOccurred())

				code1 := "name_t"
				code2 := "city_t"

				e = app.UpdateUserStringAttr(id1, code1, "andrey_t")
				Expect(e).ShouldNot(HaveOccurred())
				e = app.UpdateUserStringAttr(id1, code2, "moscow_t")
				Expect(e).ShouldNot(HaveOccurred())

				e = app.UpdateUserStringAttr(id2, code1, "andrey_t")
				Expect(e).ShouldNot(HaveOccurred())
				e = app.UpdateUserStringAttr(id2, code2, "piterburg_t")
				Expect(e).ShouldNot(HaveOccurred())
				e = app.DeleteUserStringAttr(id2, code2)
				Expect(e).ShouldNot(HaveOccurred())

				e = app.UpdateUserStringAttr(id3, code1, "andrey_t")
				Expect(e).ShouldNot(HaveOccurred())
				e = app.UpdateUserStringAttr(id3, code2, "rostov_t")
				Expect(e).ShouldNot(HaveOccurred())

				users, e := app.GetUsersByAttr(map[string]string{code1: "andrey_t"})
				Expect(e).ShouldNot(HaveOccurred())
				Expect(users).ShouldNot(BeNil())
				Expect(len(users)).Should(Equal(3))

				users, e = app.GetUsersByAttr(map[string]string{code1: "andrey_t", code2: "moscow_t"})
				Expect(e).ShouldNot(HaveOccurred())
				Expect(users).ShouldNot(BeNil())
				Expect(len(users)).Should(Equal(1))

				app.DeleteUserStringAttr(id1, code1)
				app.DeleteUserStringAttr(id1, code2)
				app.DeleteUserStringAttr(id2, code1)
				app.DeleteUserStringAttr(id2, code2)
				app.DeleteUserStringAttr(id3, code1)
				app.DeleteUserStringAttr(id3, code2)

			})
		})
	})
})
