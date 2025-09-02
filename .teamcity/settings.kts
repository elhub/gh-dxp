import jetbrains.buildServer.configs.kotlin.ArtifactRule
import no.elhub.devxp.build.configuration.pipeline.constants.Group.DEVXP
import no.elhub.devxp.build.configuration.pipeline.dsl.elhubProject
import no.elhub.devxp.build.configuration.pipeline.jobs.makeVerify

elhubProject(DEVXP, "gh-dxp") {

    params {
        param("env.PATH", "\$PATH:/usr/local/go/bin:/usr/bin")
        param("env.GOROOT", "/usr/local/go")
    }

    pipeline {
        makeVerify {
            buildArtifactRules = listOf(ArtifactRule.include("build/coverage.*", "build.zip"))
            outputArtifactRules = listOf(ArtifactRule.include("build.zip!coverage.*, build/coverage.out"))
            sonarScanSettings = {
                sonarProjectSources = "."
                additionalParams = arrayListOf("-Dsonar.go.coverage.reportPaths=build/coverage.out")
            }
            enablePublishMetrics = true
        }
    }
}
