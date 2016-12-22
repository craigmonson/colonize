package util_test

import (
	. "github.com/craigmonson/colonize/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// NOTE: this is actually testing the 'test' directory in the root of the
// project.
var _ = Describe("util/paths", func() {
	var err error
	path := "../test/"

	Describe("FindCfgPath", func() {
		Context("given a path that contains a "+RootFiles[0], func() {
			var testPath string

			BeforeEach(func() {
				testPath, err = FindCfgPath(path)
			})

			It("should return the full path to the config file", func() {
				Ω(testPath).To(Equal(path + RootFiles[0]))
			})

			It("should not return an error", func() {
				Ω(err).ShouldNot(HaveOccurred())
			})
		})

		Context("given a path that doesn't have a "+RootFiles[0], func() {
			var testPath string

			BeforeEach(func() {
				testPath, err = FindCfgPath("/etc")
			})

			It("should return empty string.", func() {
				Ω(testPath).To(Equal(""))
			})

			It("should return an error", func() {
				Ω(err).Should(MatchError(RootFiles[0] + " not found in the directory tree."))
			})
		})
	})

	Describe("GetBasename", func() {
		Context("given a path with a trailing slash", func() {
			testPath := "../test/vpc/"
			base := GetBasename(testPath)

			It("should return the proper basename", func() {
				Ω(base).To(Equal("vpc"))
			})
		})

		Context("given a path with no trailing slash", func() {
			testPath := "../test/vpc"
			base := GetBasename(testPath)

			It("should return the proper basename", func() {
				Ω(base).To(Equal("vpc"))
			})
		})
	})

	Describe("GetTmplRelPath", func() {
		Context("given a full path and a root path", func() {
			full := "../test/microservices/api"
			root := "../test"
			tmpl := GetTmplRelPath(full, root)

			It("should return the difference between the two", func() {
				Ω(tmpl).To(Equal("microservices/api"))
			})
		})
	})

	Describe("GetDir", func() {
		It("should return just the underlying directory", func() {
			p := GetDir("foo/bar/baz/buz.config")
			Ω(p).To(Equal("foo/bar/baz"))
		})
	})

	Describe("GetTreePaths", func() {
		It("should return a slice of paths", func() {
			p := GetTreePaths("foo/bar/baz")
			Ω(p).To(Equal([]string{"foo", "foo/bar", "foo/bar/baz"}))
		})
	})

	Describe("AppendPathToPaths", func() {
		It("should return a slice of paths with the filename added", func() {
			p := AppendPathToPaths([]string{"foo", "foo/bar"}, "baz")
			Ω(p).To(Equal([]string{"foo/baz", "foo/bar/baz"}))
		})
	})

	Describe("PrependPathToPaths", func() {
		It("should return a slice of paths with the path prepended", func() {
			p := PrependPathToPaths([]string{"foo", "foo/bar"}, "baz")
			Ω(p).To(Equal([]string{"baz/foo", "baz/foo/bar"}))
		})
	})

	Describe("PathJoin", func() {
		It("should join paths together", func() {
			Ω(PathJoin("foo", "bar/baz")).To(Equal("foo/bar/baz"))
		})
	})

	Describe("AddFileToWalkablePath", func() {
		It("should add file to each of the paths under the walkable path", func() {
			wSlice := AddFileToWalkablePath("foo/bar", "baz.conf")
			Ω(wSlice).To(Equal([]string{"foo/baz.conf", "foo/bar/baz.conf"}))
		})
	})
})
