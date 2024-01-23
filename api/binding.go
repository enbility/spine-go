package api

type BindingEntry struct {
	Id            uint64
	ServerFeature FeatureLocalInterface
	ClientFeature FeatureRemoteInterface
}
