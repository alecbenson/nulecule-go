package nulecule

import "testing"

func TestIsExternal(t *testing.T) {
	c := Component{}
	if ext := IsExternal(c); !ext {
		t.Fatalf("Test component should have returned isExtenal=%t, got %t", true, ext)
	}

	artifactEntry := ArtifactEntry{Path: "testPath"}
	artifactEntry2 := ArtifactEntry{Path: "testPath"}
	artifactMap := make(map[string][]ArtifactEntry)
	artifactMap["test"] = []ArtifactEntry{artifactEntry, artifactEntry2}
	c.Artifacts = artifactMap
	if ext := IsExternal(c); ext {
		t.Fatalf("Test component should have returned isExtenal=%t, got %t", false, ext)
	}
}

func TestGetSourceImage(t *testing.T) {
	c := Component{Source: "docker://expectedResult"}
	result, err := GetSourceImage(c)
	if err != nil {
		t.Fatalf("Error getting source image from component: %s", err)
	}
	expected := "expectedResult"
	if result != expected {
		t.Fatalf("Expected source image to be %s, got %s instead", expected, result)
	}

	//This test should throw an error because we should not be
	//geting the source image of a local resource
	artifactEntry := ArtifactEntry{Path: "testPath"}
	artifactMap := make(map[string][]ArtifactEntry)
	artifactMap["test"] = []ArtifactEntry{artifactEntry}
	c.Artifacts = artifactMap
	c.Source = ""
	if _, err := GetSourceImage(c); err == nil {
		t.Fatalf("Should not be able to get source image of local component")
	}

}

/*
//GetSourceImage fetches the sanitized source path of the image
func GetSourceImage(component Component) (string, error) {
	source := component.Source
	if !IsExternal(component) {
		logrus.Errorf("Cannot get external source of local component\n")
		return "", errors.New("Cannot get source of local component")
	}

	if strings.HasPrefix(source, utils.DOCKER_PREFIX) {
		return strings.TrimPrefix(source, utils.DOCKER_PREFIX), nil
	}

	logrus.Errorf("Could not get source image from component source: %v\n", component)
	return "", errors.New("Could not get source image")
}

*/
