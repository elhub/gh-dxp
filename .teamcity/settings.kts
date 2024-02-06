import no.elhub.devxp.build.configuration.pipeline.ElhubProject.Companion.elhubProject
import no.elhub.devxp.build.configuration.pipeline.constants.Group.DEVXP
import no.elhub.devxp.build.configuration.pipeline.jobs.goVerify

elhubProject(DEVXP, "gh-devxp") {

    pipeline(withReleaseVersion = false) {
        goVerify()
    }
}
