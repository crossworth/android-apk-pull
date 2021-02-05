package main

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	listPackagesCmd := exec.Command("adb", strings.Split("shell pm list packages", " ")...)
	output, err := listPackagesCmd.CombinedOutput()
	if err != nil {
		log.Fatalf("error reading package list: %s", err)
	}

	if !strings.HasPrefix(string(output), "package:") {
		log.Fatalf("got wrong output while reading packages: %s", string(output))
	}

	packages := strings.Split(removeCarriageReturn(string(output)), "\n")
	for _, pkg := range packages {
		pkgName := strings.ReplaceAll(pkg, "package:", "")
		log.Printf("pulling %s\n", pkgName)
		extractPackage(pkgName)
	}
}

func extractPackage(pkgName string) {
	log.Printf("resolving package path from %s\n", pkgName)
	packagePathCmd := exec.Command("adb", strings.Split("shell pm path "+pkgName, " ")...)
	output, err := packagePathCmd.CombinedOutput()
	if err != nil {
		log.Printf("error reading package path from %s: %s", pkgName, err)
		return
	}

	if !strings.HasPrefix(string(output), "package:") {
		log.Printf("got wrong output while extracting package %s: %s", pkgName, string(output))
		return
	}

	// package:/data/app/com.amazon.kindle-Y9ANE6dYYuEx6u8mIIDoYQ==/base.apk
	pkgPath := strings.ReplaceAll(removeLineFeed(removeCarriageReturn(string(output))), "package:", "")

	log.Printf("pulling %s\n", pkgPath)
	pullCmd := exec.Command("adb", strings.Split("pull "+pkgPath, " ")...)
	output, err = pullCmd.CombinedOutput()
	if err != nil {
		log.Printf("error pulling package %s: %s", pkgName, err)
		return
	}

	_ = os.Rename("base.apk", pkgName+".apk")

	log.Printf("pulled %s\n", pkgName)
}

func removeCarriageReturn(input string) string {
	return strings.ReplaceAll(input, "\r", "")
}

func removeLineFeed(input string) string {
	return strings.ReplaceAll(input, "\n", "")
}
