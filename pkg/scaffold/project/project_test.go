package project

import (
	"path/filepath"

	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/tsungming/controller-tools/pkg/scaffold"
	"github.com/tsungming/controller-tools/pkg/scaffold/input"
	"github.com/tsungming/controller-tools/pkg/scaffold/scaffoldtest"
)

var _ = Describe("Project", func() {
	var result *scaffoldtest.TestResult
	var writeToPath, goldenPath string
	var s *scaffold.Scaffold

	JustBeforeEach(func() {
		s, result = scaffoldtest.NewTestScaffold(writeToPath, goldenPath)
	})

	Describe("scaffolding a boilerplate file", func() {
		BeforeEach(func() {
			goldenPath = filepath.Join("hack", "boilerplate.go.txt")
			writeToPath = goldenPath
		})

		It("should match the golden file", func() {
			instance := &Boilerplate{Year: "2018", License: "apache2", Owner: "The Kubernetes authors"}
			Expect(s.Execute(input.Options{}, instance)).NotTo(HaveOccurred())
			Expect(result.Actual.String()).To(BeEquivalentTo(result.Golden))
		})

		It("should skip writing boilerplate if the file exists", func() {
			i, err := (&Boilerplate{}).GetInput()
			Expect(err).NotTo(HaveOccurred())
			Expect(i.IfExistsAction).To(Equal(input.Skip))
		})

		Context("for apache2", func() {
			It("should write the apache2 boilerplate with specified owners", func() {
				instance := &Boilerplate{Year: "2018", Owner: "Example Owners"}
				Expect(s.Execute(input.Options{}, instance)).NotTo(HaveOccurred())
				e := strings.Replace(
					result.Golden, "The Kubernetes authors", "Example Owners", -1)
				Expect(result.Actual.String()).To(BeEquivalentTo(e))
			})

			It("should use apache2 as the default", func() {
				instance := &Boilerplate{Year: "2018", Owner: "The Kubernetes authors"}
				Expect(s.Execute(input.Options{}, instance)).NotTo(HaveOccurred())
				Expect(result.Actual.String()).To(BeEquivalentTo(result.Golden))
			})
		})

		Context("for none", func() {
			It("should write the empty boilerplate", func() {
				// Scaffold a boilerplate file
				instance := &Boilerplate{Year: "2019", License: "none", Owner: "Example Owners"}
				Expect(s.Execute(input.Options{}, instance)).NotTo(HaveOccurred())
				Expect(result.Actual.String()).To(BeEquivalentTo(`/*
Copyright 2019 Example Owners.
*/`))
			})
		})

		Context("if the boilerplate is given", func() {
			It("should skip writing Gopkg.toml", func() {
				instance := &Boilerplate{}
				instance.Boilerplate = `/* Hello World */`

				Expect(s.Execute(input.Options{}, instance)).NotTo(HaveOccurred())
				Expect(result.Actual.String()).To(BeEquivalentTo(`/* Hello World */`))
			})
		})
	})
})
