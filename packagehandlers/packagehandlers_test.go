package packagehandlers

import (
	"fmt"
	testdatautils "github.com/jfrog/build-info-go/build/testdata"
	"github.com/jfrog/frogbot/utils"
	"github.com/jfrog/jfrog-cli-core/v2/utils/coreutils"
	"github.com/jfrog/jfrog-cli-core/v2/xray/formats"
	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

type dependencyFixTest struct {
	vulnDetails           *utils.VulnerabilityDetails
	scanDetails           *utils.ScanDetails
	fixSupported          bool
	specificTechVersion   string
	uniqueChecksExtraArgs []string
}

const (
	requirementsFile    = "oslo.config>=1.12.1,<1.13\noslo.utils<5.0,>=4.0.0\nparamiko==2.7.2\npasslib<=1.7.4\nprance>=0.9.0\nprompt-toolkit~=1.0.15\npyinotify>0.9.6\nPyJWT>1.7.1\nurllib3 > 1.1.9, < 1.5.*"
	GoPackageDescriptor = "go.mod"
)

type pipPackageRegexTest struct {
	packageName         string
	expectedRequirement string
}

func TestUpdateDependency(t *testing.T) {
	testCases := [][]dependencyFixTest{
		// Go test cases
		{
			{
				vulnDetails: &utils.VulnerabilityDetails{
					SuggestedFixedVersion:       "0.0.0-20201216223049-8b5274cf687f",
					IsDirectDependency:          false,
					VulnerabilityOrViolationRow: formats.VulnerabilityOrViolationRow{Technology: coreutils.Go, ImpactedDependencyName: "golang.org/x/crypto"},
				},
				fixSupported:          true,
				uniqueChecksExtraArgs: []string{GoPackageDescriptor},
			},
			{
				vulnDetails: &utils.VulnerabilityDetails{
					SuggestedFixedVersion:       "1.7.7",
					IsDirectDependency:          true,
					VulnerabilityOrViolationRow: formats.VulnerabilityOrViolationRow{Technology: coreutils.Go, ImpactedDependencyName: "github.com/gin-gonic/gin"},
				},
				fixSupported:          true,
				uniqueChecksExtraArgs: []string{GoPackageDescriptor},
			},
			{
				vulnDetails: &utils.VulnerabilityDetails{
					SuggestedFixedVersion:       "1.3.0",
					IsDirectDependency:          true,
					VulnerabilityOrViolationRow: formats.VulnerabilityOrViolationRow{Technology: coreutils.Go, ImpactedDependencyName: "github.com/google/uuid"},
				},
				fixSupported:          true,
				uniqueChecksExtraArgs: []string{GoPackageDescriptor},
			},
		},

		// Python test cases (includes pip, pipenv, poetry)
		{
			{
				vulnDetails: &utils.VulnerabilityDetails{
					SuggestedFixedVersion:       "1.25.9",
					VulnerabilityOrViolationRow: formats.VulnerabilityOrViolationRow{Technology: coreutils.Pip, ImpactedDependencyName: "urllib3"},
				},
				scanDetails:  &utils.ScanDetails{Project: &utils.Project{PipRequirementsFile: "requirements.txt"}},
				fixSupported: false,
			},
			{
				vulnDetails: &utils.VulnerabilityDetails{
					SuggestedFixedVersion:       "1.25.9",
					VulnerabilityOrViolationRow: formats.VulnerabilityOrViolationRow{Technology: coreutils.Poetry, ImpactedDependencyName: "urllib3"},
				},
				scanDetails:  &utils.ScanDetails{Project: &utils.Project{PipRequirementsFile: "pyproejct.toml"}},
				fixSupported: false,
			},
			{
				vulnDetails: &utils.VulnerabilityDetails{
					SuggestedFixedVersion:       "1.25.9",
					VulnerabilityOrViolationRow: formats.VulnerabilityOrViolationRow{Technology: coreutils.Pipenv, ImpactedDependencyName: "urllib3"},
				},
				scanDetails:  &utils.ScanDetails{Project: &utils.Project{PipRequirementsFile: "Pipfile"}},
				fixSupported: false,
			},
			{
				vulnDetails: &utils.VulnerabilityDetails{
					SuggestedFixedVersion:       "2.4.0",
					VulnerabilityOrViolationRow: formats.VulnerabilityOrViolationRow{Technology: coreutils.Pip, ImpactedDependencyName: "pyjwt"},
					IsDirectDependency:          true,
				},
				scanDetails:  &utils.ScanDetails{Project: &utils.Project{PipRequirementsFile: "requirements.txt"}},
				fixSupported: true,
			},
			{
				vulnDetails: &utils.VulnerabilityDetails{
					SuggestedFixedVersion:       "2.4.0",
					VulnerabilityOrViolationRow: formats.VulnerabilityOrViolationRow{Technology: coreutils.Pip, ImpactedDependencyName: "Pyjwt"},
					IsDirectDependency:          true},
				scanDetails:  &utils.ScanDetails{Project: &utils.Project{PipRequirementsFile: "requirements.txt"}},
				fixSupported: true,
			},
			{
				vulnDetails: &utils.VulnerabilityDetails{
					SuggestedFixedVersion:       "2.4.0",
					VulnerabilityOrViolationRow: formats.VulnerabilityOrViolationRow{Technology: coreutils.Pip, ImpactedDependencyName: "pyjwt"},
					IsDirectDependency:          true},
				scanDetails:  &utils.ScanDetails{Project: &utils.Project{PipRequirementsFile: "setup.py"}},
				fixSupported: true,
			},
			{
				vulnDetails: &utils.VulnerabilityDetails{
					SuggestedFixedVersion:       "2.4.0",
					VulnerabilityOrViolationRow: formats.VulnerabilityOrViolationRow{Technology: coreutils.Poetry, ImpactedDependencyName: "pyjwt"},
					IsDirectDependency:          true},
				scanDetails:  &utils.ScanDetails{Project: &utils.Project{PipRequirementsFile: "pyproject.toml"}},
				fixSupported: true,
			},
		},

		// Npm test cases
		{
			{
				vulnDetails: &utils.VulnerabilityDetails{
					SuggestedFixedVersion:       "0.8.4",
					VulnerabilityOrViolationRow: formats.VulnerabilityOrViolationRow{Technology: coreutils.Npm, ImpactedDependencyName: "mpath"},
				},
				fixSupported: false,
			},
			{
				vulnDetails: &utils.VulnerabilityDetails{
					SuggestedFixedVersion:       "3.0.2",
					IsDirectDependency:          true,
					VulnerabilityOrViolationRow: formats.VulnerabilityOrViolationRow{Technology: coreutils.Npm, ImpactedDependencyName: "minimatch"},
				},
				fixSupported: true,
			},
		},

		// Yarn test cases
		{
			{
				// This test case directs to non-existing directory. It only checks if the dependency update is blocked if the vulnerable dependency is not a direct dependency
				vulnDetails: &utils.VulnerabilityDetails{
					SuggestedFixedVersion:       "1.2.6",
					IsDirectDependency:          false,
					VulnerabilityOrViolationRow: formats.VulnerabilityOrViolationRow{Technology: coreutils.Yarn, ImpactedDependencyName: "minimist"},
				},
				fixSupported: false,
			},
			{
				vulnDetails: &utils.VulnerabilityDetails{
					SuggestedFixedVersion:       "1.2.6",
					IsDirectDependency:          true,
					VulnerabilityOrViolationRow: formats.VulnerabilityOrViolationRow{Technology: coreutils.Yarn, ImpactedDependencyName: "minimist"},
				},
				fixSupported:        true,
				specificTechVersion: "1",
			},
			{
				vulnDetails: &utils.VulnerabilityDetails{
					SuggestedFixedVersion:       "1.2.6",
					IsDirectDependency:          true,
					VulnerabilityOrViolationRow: formats.VulnerabilityOrViolationRow{Technology: coreutils.Yarn, ImpactedDependencyName: "minimist"},
				},
				fixSupported:        true,
				specificTechVersion: "2",
			},
		},

		// Maven test cases
		{
			{
				vulnDetails: &utils.VulnerabilityDetails{
					SuggestedFixedVersion:       "2.7",
					VulnerabilityOrViolationRow: formats.VulnerabilityOrViolationRow{Technology: coreutils.Maven, ImpactedDependencyName: "commons-io:commons-io"},
					IsDirectDependency:          true},
				scanDetails:  &utils.ScanDetails{Project: &utils.Project{DepsRepo: ""}, ServerDetails: nil},
				fixSupported: true,
			},
			{
				vulnDetails: &utils.VulnerabilityDetails{
					SuggestedFixedVersion:       "4.3.20",
					VulnerabilityOrViolationRow: formats.VulnerabilityOrViolationRow{Technology: coreutils.Maven, ImpactedDependencyName: "org.springframework:spring-core"},
					IsDirectDependency:          false},
				scanDetails:  &utils.ScanDetails{Project: &utils.Project{DepsRepo: ""}, ServerDetails: nil},
				fixSupported: false,
			},
		},

		// NuGet test cases
		{
			{
				// This test case directs to non-existing directory. It only checks if the dependency update is blocked if the vulnerable dependency is not a direct dependency
				vulnDetails: &utils.VulnerabilityDetails{
					SuggestedFixedVersion:       "1.1.1",
					IsDirectDependency:          false,
					VulnerabilityOrViolationRow: formats.VulnerabilityOrViolationRow{Technology: coreutils.Nuget, ImpactedDependencyName: "snappier"},
				},
				fixSupported: false,
			},
			{
				vulnDetails: &utils.VulnerabilityDetails{
					SuggestedFixedVersion:       "1.1.1",
					IsDirectDependency:          true,
					VulnerabilityOrViolationRow: formats.VulnerabilityOrViolationRow{Technology: coreutils.Nuget, ImpactedDependencyName: "snappier"},
				},
				fixSupported: true,
			},
		},
	}

	for _, testBatch := range testCases {
		for _, test := range testBatch {
			packageHandler := GetCompatiblePackageHandler(test.vulnDetails, test.scanDetails)
			t.Run(fmt.Sprintf("%s:%s direct:%s", test.vulnDetails.Technology.ToString()+test.specificTechVersion, test.vulnDetails.ImpactedDependencyName, strconv.FormatBool(test.vulnDetails.IsDirectDependency)),
				func(t *testing.T) {
					testDataDir := getTestDataDir(t, test.vulnDetails.IsDirectDependency)
					cleanup := createTempDirAndChdir(t, testDataDir, test.vulnDetails.Technology.ToString()+test.specificTechVersion)
					defer cleanup()
					err := packageHandler.UpdateDependency(test.vulnDetails)
					if test.fixSupported {
						assert.NoError(t, err)
						uniquePackageManagerChecks(t, test)
					} else {
						assert.Error(t, err)
						assert.IsType(t, &utils.ErrUnsupportedFix{}, err, "Expected unsupported fix error")
					}
				})
		}
	}

}

func TestPipPackageRegex(t *testing.T) {
	var pipPackagesRegexTests = []pipPackageRegexTest{
		{"oslo.config", "oslo.config>=1.12.1,<1.13"},
		{"oslo.utils", "oslo.utils<5.0,>=4.0.0"},
		{"paramiko", "paramiko==2.7.2"},
		{"passlib", "passlib<=1.7.4"},
		{"PassLib", "passlib<=1.7.4"},
		{"prance", "prance>=0.9.0"},
		{"prompt-toolkit", "prompt-toolkit~=1.0.15"},
		{"pyinotify", "pyinotify>0.9.6"},
		{"pyjwt", "pyjwt>1.7.1"},
		{"PyJWT", "pyjwt>1.7.1"},
		{"urllib3", "urllib3 > 1.1.9, < 1.5.*"},
	}
	for _, pack := range pipPackagesRegexTests {
		re := regexp.MustCompile(PythonPackageRegexPrefix + "(" + pack.packageName + "|" + strings.ToLower(pack.packageName) + ")" + PythonPackageRegexSuffix)
		found := re.FindString(requirementsFile)
		assert.Equal(t, pack.expectedRequirement, strings.ToLower(found))
	}
}

// Maven utils functions
func TestGetDependenciesFromPomXmlSingleDependency(t *testing.T) {
	testCases := []string{`<dependency>
	<groupId>org.apache.commons</groupId>
	<artifactId>commons-email</artifactId>
	<version>1.1</version>
	<scope>compile</scope>
</dependency>`,
		`<dependency>
	<groupId> 	org.apache.commons</groupId>
	<artifactId>commons-email	 </artifactId>
	<version>  1.1  </version>
	<scope> compile		</scope>
</dependency>`,
	}

	for _, testCase := range testCases {
		result, err := getMavenDependencies([]byte(testCase))
		assert.NoError(t, err)

		assert.Len(t, result, 1)
		assert.Equal(t, "org.apache.commons", result[0].GroupId)
		assert.Equal(t, "commons-email", result[0].ArtifactId)
		assert.Equal(t, "1.1", result[0].Version)
	}
}

func TestGetDependenciesFromPomXmlMultiDependency(t *testing.T) {
	testCases := []string{`
<project xmlns="http://maven.apache.org/POM/4.0.0" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
        xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 https://maven.apache.org/maven-v4_0_0.xsd">
<dependencies>
       <dependency>
           <groupId>org.apache.commons</groupId>
           <artifactId>commons-email</artifactId>
           <version>1.1</version>
           <scope>compile</scope>
       </dependency>
       <dependency>
           <groupId>org.codehaus.plexus</groupId>
           <artifactId>plexus-utils</artifactId>
           <version>1.5.1</version>
       </dependency>
	</dependencies>
</project>`,
	}

	for _, testCase := range testCases {
		result, err := getMavenDependencies([]byte(testCase))
		assert.NoError(t, err)

		assert.Len(t, result, 2)
		assert.Equal(t, "org.apache.commons", result[0].GroupId)
		assert.Equal(t, "commons-email", result[0].ArtifactId)
		assert.Equal(t, "1.1", result[0].Version)

		assert.Equal(t, "org.codehaus.plexus", result[1].GroupId)
		assert.Equal(t, "plexus-utils", result[1].ArtifactId)
		assert.Equal(t, "1.5.1", result[1].Version)
	}
}

func TestGetPluginsFromPomXml(t *testing.T) {
	testCase :=
		`<project>
			<build>
       <plugins>
           <plugin>
               <groupId>org.apache.maven.plugins</groupId>
               <artifactId>maven-source-plugin</artifactId>
           </plugin>
           <plugin>
               <groupId>com.github.spotbugs</groupId>
               <artifactId>spotbugs-maven-plugin</artifactId>
               <version>4.5.3.0</version>
               <configuration>
                   <excludeFilterFile>spotbugs-security-exclude.xml</excludeFilterFile>
                   <plugins>
                       <plugin>
                           <groupId>com.h3xstream.findsecbugs</groupId>
                           <artifactId>findsecbugs-plugin</artifactId>
                           <version>1.12.0</version>
                       </plugin>
                   </plugins>
               </configuration>
           </plugin>
           <plugin>
               <groupId>org.apache.maven.plugins</groupId>
               <artifactId>maven-surefire-plugin</artifactId>
               <version>2.22.1</version>
               <configuration>
                   <systemPropertyVariables>
                       <!--This will disable JenkinsRule timeout-->
                       <maven.surefire.debug>true</maven.surefire.debug>
                   </systemPropertyVariables>
                   <excludes>
                       <exclude>**/InjectedTest.java</exclude>
                       <exclude>**/*ITest.java</exclude>
                   </excludes>
               </configuration>
           </plugin>
       </plugins>
   </build>
	</project>
		`
	plugins, err := getMavenDependencies([]byte(testCase))
	assert.NoError(t, err)
	assert.Equal(t, "org.apache.maven.plugins", plugins[0].GroupId)
	assert.Equal(t, "maven-source-plugin", plugins[0].ArtifactId)
	assert.Equal(t, "com.github.spotbugs", plugins[1].GroupId)
	assert.Equal(t, "spotbugs-maven-plugin", plugins[1].ArtifactId)
	assert.Equal(t, "4.5.3.0", plugins[1].Version)
	assert.Equal(t, "com.h3xstream.findsecbugs", plugins[2].GroupId)
	assert.Equal(t, "findsecbugs-plugin", plugins[2].ArtifactId)
	assert.Equal(t, "1.12.0", plugins[2].Version)
	assert.Equal(t, "org.apache.maven.plugins", plugins[3].GroupId)
	assert.Equal(t, "maven-surefire-plugin", plugins[3].ArtifactId)
	assert.Equal(t, "2.22.1", plugins[3].Version)
}

func TestGetDependenciesFromDependencyManagement(t *testing.T) {
	testCase := `
<project xmlns="http://maven.apache.org/POM/4.0.0" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
        xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 https://maven.apache.org/maven-v4_0_0.xsd">
   <dependencyManagement>
       <dependencies>
           <dependency>
               <groupId>io.jenkins.tools.bom</groupId>
               <artifactId>bom-2.346.x</artifactId>
               <version>1607.va_c1576527071</version>
               <scope>import</scope>
               <type>pom</type>
           </dependency>
           <dependency>
               <groupId>com.fasterxml.jackson.core</groupId>
               <artifactId>jackson-core</artifactId>
               <version>2.13.4</version>
           </dependency>
           <dependency>
               <groupId>com.fasterxml.jackson.core</groupId>
               <artifactId>jackson-databind</artifactId>
               <version>2.13.4.2</version>
           </dependency>
           <dependency>
               <groupId>com.fasterxml.jackson.core</groupId>
               <artifactId>jackson-annotations</artifactId>
               <version>2.13.4</version>
           </dependency>
           <dependency>
               <groupId>org.apache.httpcomponents</groupId>
               <artifactId>httpcore</artifactId>
               <version>4.4.15</version>
           </dependency>
           <dependency>
               <groupId>org.jenkins-ci.plugins.workflow</groupId>
               <artifactId>workflow-durable-task-step</artifactId>
               <version>1190.vc93d7d457042</version>
               <scope>test</scope>
           </dependency>
       </dependencies>
   </dependencyManagement>
</project>
`
	dependencies, err := getMavenDependencies([]byte(testCase))
	assert.NoError(t, err)
	assert.Len(t, dependencies, 6)
	for _, dependency := range dependencies {
		assert.True(t, dependency.foundInDependencyManagement)
	}
}

func TestMavenGavReader(t *testing.T) {
	mvnHandler := &MavenPackageHandler{}
	currDir, err := os.Getwd()
	assert.NoError(t, err)
	tmpDir, err := os.MkdirTemp("", "")
	assert.NoError(t, err)
	assert.NoError(t, fileutils.CopyDir(filepath.Join("..", "testdata", "projects", "maven"), tmpDir, true, nil))
	assert.NoError(t, os.Chdir(tmpDir))
	defer func() {
		assert.NoError(t, os.Chdir(currDir))
	}()
	// Test installMavenGavReader
	assert.NoError(t, mvnHandler.installMavenGavReader())
	assert.True(t, mvnHandler.isMavenGavReaderInstalled)
	// Test getProjectPoms using the maven-gav-reader plugin
	assert.NoError(t, mvnHandler.getProjectPoms())
	assert.Len(t, mvnHandler.pomPaths, 2)
}

// General Utils functions
func TestFixVersionInfo_UpdateFixVersionIfMax(t *testing.T) {
	type testCase struct {
		fixVersionInfo utils.VulnerabilityDetails
		newFixVersion  string
		expectedOutput string
	}

	testCases := []testCase{
		{fixVersionInfo: utils.VulnerabilityDetails{SuggestedFixedVersion: "1.2.3", IsDirectDependency: true}, newFixVersion: "1.2.4", expectedOutput: "1.2.4"},
		{fixVersionInfo: utils.VulnerabilityDetails{SuggestedFixedVersion: "1.2.3", IsDirectDependency: true}, newFixVersion: "1.0.4", expectedOutput: "1.2.3"},
	}

	for _, tc := range testCases {
		t.Run(tc.expectedOutput, func(t *testing.T) {
			tc.fixVersionInfo.UpdateFixVersionIfMax(tc.newFixVersion)
			assert.Equal(t, tc.expectedOutput, tc.fixVersionInfo.SuggestedFixedVersion)
		})
	}
}

func TestUpdatePackageVersion(t *testing.T) {
	testProjectPath := filepath.Join("..", "testdata", "packagehandlers")
	currDir, err := os.Getwd()
	assert.NoError(t, err)
	tmpDir, err := os.MkdirTemp("", "")
	assert.NoError(t, err)
	assert.NoError(t, fileutils.CopyDir(testProjectPath, tmpDir, true, nil))
	assert.NoError(t, os.Chdir(tmpDir))
	defer func() {
		assert.NoError(t, os.Chdir(currDir))
	}()
	testCases := []struct {
		impactedPackage             string
		fixedVersion                string
		foundInDependencyManagement bool
	}{
		{impactedPackage: "org.jfrog.filespecs:file-specs-java", fixedVersion: "1.1.2"},
		{impactedPackage: "com.fasterxml.jackson.core:jackson-core", fixedVersion: "2.15.0", foundInDependencyManagement: true},
		{impactedPackage: "org.apache.httpcomponents:httpcore", fixedVersion: "4.4.16", foundInDependencyManagement: true},
	}
	mvnHandler := &MavenPackageHandler{}
	for _, test := range testCases {
		assert.NoError(t, mvnHandler.updatePackageVersion(test.impactedPackage, test.fixedVersion, test.foundInDependencyManagement))
	}
	modifiedPom, err := os.ReadFile("pom.xml")
	assert.NoError(t, err)
	for _, test := range testCases {
		assert.Contains(t, fmt.Sprintf("<version>%s</version>", string(modifiedPom)), test.fixedVersion)
	}
}

func TestUpdatePropertiesVersion(t *testing.T) {
	testProjectPath := filepath.Join("..", "testdata", "packagehandlers")
	currDir, err := os.Getwd()
	assert.NoError(t, err)
	tmpDir, err := os.MkdirTemp("", "")
	assert.NoError(t, err)
	assert.NoError(t, fileutils.CopyDir(testProjectPath, tmpDir, true, nil))
	assert.NoError(t, os.Chdir(tmpDir))
	defer func() {
		assert.NoError(t, os.Chdir(currDir))
	}()
	mvnHandler := &MavenPackageHandler{}
	assert.NoError(t, mvnHandler.updateProperties(&pomDependencyDetails{properties: []string{"buildinfo.version"}}, "2.39.9"))
	modifiedPom, err := os.ReadFile("pom.xml")
	assert.NoError(t, err)
	assert.Contains(t, string(modifiedPom), "2.39.9")
}

func getTestDataDir(t *testing.T, directDependency bool) string {
	var projectDir string
	if directDependency {
		projectDir = "projects"
	} else {
		projectDir = "indirect-projects"
	}
	testdataDir, err := filepath.Abs(filepath.Join("..", "testdata", projectDir))
	assert.NoError(t, err)
	return testdataDir
}

func createTempDirAndChdir(t *testing.T, testdataDir string, tech string) func() {
	// Create temp technology project
	projectPath := filepath.Join(testdataDir, tech)
	tmpProjectPath, cleanup := testdatautils.CreateTestProject(t, projectPath)
	currDir, err := os.Getwd()
	assert.NoError(t, err)
	assert.NoError(t, os.Chdir(tmpProjectPath))
	return func() {
		cleanup()
		assert.NoError(t, os.Chdir(currDir))
	}
}

func assertFixVersionInPackageDescriptor(t *testing.T, test dependencyFixTest, packageDescriptor string) {
	file, err := os.ReadFile(packageDescriptor)
	assert.NoError(t, err)
	if test.fixSupported {
		assert.Contains(t, string(file), test.vulnDetails.SuggestedFixedVersion)
		// Verify that case-sensitive packages in python are lowered
		assert.Contains(t, string(file), strings.ToLower(test.vulnDetails.ImpactedDependencyName))
	} else {
		assert.NotContains(t, string(file), test.vulnDetails)
	}
}

// This function is intended to add unique checks for specific package managers
func uniquePackageManagerChecks(t *testing.T, test dependencyFixTest) {
	technology := test.vulnDetails.Technology
	extraArgs := test.uniqueChecksExtraArgs
	switch technology {
	case coreutils.Go:
		packageDescriptor := extraArgs[0]
		assertFixVersionInPackageDescriptor(t, test, packageDescriptor)
	default:
	}
}

func TestGetFixedPackage(t *testing.T) {
	var testcases = []struct {
		impactedPackage       string
		versionOperator       string
		suggestedFixedVersion string
		expectedOutput        []string
	}{
		{
			impactedPackage:       "snappier",
			versionOperator:       " -v ",
			suggestedFixedVersion: "1.1.1",
			expectedOutput:        []string{"snappier", "-v", "1.1.1"},
		},
		{
			impactedPackage:       "json",
			versionOperator:       "@",
			suggestedFixedVersion: "10.0.0",
			expectedOutput:        []string{"json@10.0.0"},
		},
	}

	for _, test := range testcases {
		fixedPackageArgs := getFixedPackage(test.impactedPackage, test.versionOperator, test.suggestedFixedVersion)
		assert.Equal(t, test.expectedOutput, fixedPackageArgs)
	}
}
