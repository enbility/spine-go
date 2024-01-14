package api

type BindingEntry struct {
	Id            uint64
	ServerFeature FeatureLocal
	ClientFeature FeatureRemote
}
