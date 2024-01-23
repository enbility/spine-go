package api

type SubscriptionEntry struct {
	Id            uint64
	ServerFeature FeatureLocalInterface
	ClientFeature FeatureRemoteInterface
}
