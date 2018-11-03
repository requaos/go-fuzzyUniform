package fuzzy

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCompute(t *testing.T) {
	hash, err := Compute([]byte("This string might be long enough, maybe now it could be long enough"), 5)
	require.NoError(t, err)
	t.Log(hash)
	hash2, err := Compute([]byte("This string might not be long enough, maybe now it could be long enough"), 5)
	require.NoError(t, err)
	t.Log(hash2)
	result, err := Compare(hash, hash2)
	require.NoError(t, err)
	t.Log(result)
}

func TestCompare(t *testing.T) {
	hash, err := Compute([]byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit. Cras fringilla eleifend leo, quis venenatis lectus tincidunt vel. Quisque quis nibh sed enim imperdiet maximus eu ut odio. Etiam sed auctor tortor, id vestibulum nisi. Aliquam vulputate sem in dolor pharetra, in tempor mi faucibus. Aliquam ultricies, sapien eu faucibus luctus, dolor nunc porttitor enim, ac rhoncus massa tortor non tortor. Mauris sem risus, eleifend ac iaculis at, condimentum quis metus. Morbi egestas tellus nulla, sed vestibulum mauris efficitur nec. Ut consectetur tincidunt eros eget pharetra. Mauris eu ante non ante convallis tempor in sed mi. Nullam id eros sapien. Curabitur a quam condimentum, finibus enim sit amet, egestas dolor."), 15)
	require.NoError(t, err)
	t.Log(hash)
	hash2, err := Compute([]byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit. Cras fringilla eleifend leo, quis venenatis lectus tincidunt vel. Quisque quis nibh sed enim imperdiet maximus eu ut odio. Etiam sed auctor tortor, id vestibulum nisi. Aliquam vulputate sem in dolor pharetra, in tempor mi faucibus. Aliquam ultricies, sapien eu faucibus luctus, dolor nunc porttitor enim, ac rhoncus massa tortor non tortor. Mauris sem risus, eleifend ac iaculis at, condimentum quis metus. Morbi egestas tellus nulla, sed vestibulum mauris efficitur nec. Ut consectetur tincidunt eros eget pharetra. Mauris ante non ante convallis tempor in sed mi. Nullam id eros sapien. Curabitur a quam condimentum, finibus enim sit amet, egestas dolor."), 15)
	require.NoError(t, err)
	t.Log(hash2)
	result, err := Compare(hash, hash2)
	require.NoError(t, err)
	t.Log(result)
}
