package pandoc

type config struct {
	IsFullTookit bool
}

const BitbucketURLEval = "https://bitbucket.org/hwuligans/gputeachingkit-labs/src/master/"
const BitbucketURLFull = "https://bitbucket.org/hwuligans/gputeachingkit-labs-full/src/master/"

var Config = config{
	IsFullTookit: true,
}
