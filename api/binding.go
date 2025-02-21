package api

type BindingEntry struct {
	Id            uint64
	LocalFeature  FeatureLocalInterface
	RemoteFeature FeatureRemoteInterface
}
