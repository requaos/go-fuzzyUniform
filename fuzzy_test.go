package fuzzy

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCompute(t *testing.T) {
	testdata := []testdata{
		testdata{
			str:  "This string might be long enough, maybe now it could be long enough",
			str1: "This string might not be long enough, maybe now it could be long enough",
		},
		testdata{
			str:  "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Cras fringilla eleifend leo, quis venenatis lectus tincidunt vel. Quisque quis nibh sed enim imperdiet maximus eu ut odio. Etiam sed auctor tortor, id vestibulum nisi. Aliquam vulputate sem in dolor pharetra, in tempor mi faucibus. Aliquam ultricies, sapien eu faucibus luctus, dolor nunc porttitor enim, ac rhoncus massa tortor non tortor. Mauris sem risus, eleifend ac iaculis at, condimentum quis metus. Morbi egestas tellus nulla, sed vestibulum mauris efficitur nec. Ut consectetur tincidunt eros eget pharetra. Mauris eu ante non ante convallis tempor in sed mi. Nullam id eros sapien. Curabitur a quam condimentum, finibus enim sit amet, egestas dolor.",
			str1: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Cras fringilla eleifend leo, quis venenatis lectus tincidunt vel. Quisque quis nibh sed enim imperdiet maximus eu ut odio. Etiam sed auctor tortor, id vestibulum nisi. Aliquam vulputate sem in dolor pharetra, in tempor mi faucibus. Aliquam ultricies, sapien eu faucibus luctus, dolor nunc porttitor enim, ac rhoncus massa tortor non tortor. Mauris sem risus, eleifend ac iaculis at, condimentum quis metus. Morbi egestas tellus nulla, sed vestibulum mauris efficitur nec. Ut consectetur tincidunt eros eget pharetra. Mauris ante non ante convallis tempor in sed mi. Nullam id eros sapien. Curabitur a quam condimentum, finibus enim sit amet, egestas dolor.",
		},
	}
	for j := range testdata {
		t.Log(testdata[j].str)
		t.Log(testdata[j].str1)
		i := 3
		for i <= testdata[j].minFactor() {
			hash, err := Compute([]byte(testdata[j].str), i)
			require.NoError(t, err)
			t.Log(hash)

			hash2, err := Compute([]byte(testdata[j].str1), i)
			require.NoError(t, err)
			t.Log(hash2)
			result, err := CompareSimilarity(hash, hash2, Maximum)
			require.NoError(t, err)
			t.Log(result)
			hashString, hash2String := hash.String(), hash2.String()
			scale := maxLength(hashString, hash2String)
			t.Log(100 - (((distance(hashString, hash2String)*scale)/(len(hashString)+len(hash2String)))*100)/scale)
			i += 2
		}
	}
}

func TestCompare(t *testing.T) {
	hash, err := Compute([]byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit. Cras fringilla eleifend leo, quis venenatis lectus tincidunt vel. Quisque quis nibh sed enim imperdiet maximus eu ut odio. Etiam sed auctor tortor, id vestibulum nisi. Aliquam vulputate sem in dolor pharetra, in tempor mi faucibus. Aliquam ultricies, sapien eu faucibus luctus, dolor nunc porttitor enim, ac rhoncus massa tortor non tortor. Mauris sem risus, eleifend ac iaculis at, condimentum quis metus. Morbi egestas tellus nulla, sed vestibulum mauris efficitur nec. Ut consectetur tincidunt eros eget pharetra. Mauris eu ante non ante convallis tempor in sed mi. Nullam id eros sapien. Curabitur a quam condimentum, finibus enim sit amet, egestas dolor."), 19)
	require.NoError(t, err)
	t.Log(hash)
	hash2, err := Compute([]byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit. Cras fringilla eleifend leo, quis venenatis lectus tincidunt vel. Quisque quis nibh sed enim imperdiet maximus eu ut odio. Etiam sed auctor tortor, id vestibulum nisi. Aliquam vulputate sem in dolor pharetra, in tempor mi faucibus. Aliquam ultricies, sapien eu faucibus luctus, dolor nunc porttitor enim, ac rhoncus massa tortor non tortor. Mauris sem risus, eleifend ac iaculis at, condimentum quis metus. Morbi egestas tellus nulla, sed vestibulum mauris efficitur nec. Ut consectetur tincidunt eros eget pharetra. Mauris ante non ante convallis tempor in sed mi. Nullam id eros sapien. Curabitur a quam condimentum, finibus enim sit amet, egestas dolor."), 19)
	require.NoError(t, err)
	t.Log(hash2)
	result, err := Compare(hash, hash2)
	require.NoError(t, err)
	t.Log(result)
	hashString, hash2String := hash.String(), hash2.String()
	test, err := ParseHash(hashString)
	require.NoError(t, err)
	test2, err := ParseHash(hash2String)
	require.NoError(t, err)
	_, err = Compare(test, test2)
	require.NoError(t, err)
	scale := maxLength(hashString, hash2String)
	t.Log(100 - (((distance(hashString, hash2String)*scale)/(len(hashString)+len(hash2String)))*100)/scale)
}

type testdata struct {
	str, str1 string
}

func (d *testdata) minFactor() int {
	if len(d.str) < len(d.str1) {
		return len(d.str)
	}
	return len(d.str1)
}

func maxLength(a, b string) int {
	if len(a) < len(b) {
		return len(b)
	}
	return len(a)
}
