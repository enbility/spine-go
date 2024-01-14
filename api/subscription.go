package api

type SubscriptionEntry struct {
	Id            uint64
	ServerFeature FeatureLocal
	ClientFeature FeatureRemote
}
