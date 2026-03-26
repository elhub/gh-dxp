import jetbrains.buildServer.configs.kotlin.ArtifactRule
import jetbrains.buildServer.configs.kotlin.buildSteps.script
import no.elhub.devxp.build.configuration.pipeline.constants.AgentScope.LinuxAgentContext
import no.elhub.devxp.build.configuration.pipeline.constants.Group.DEVXP
import no.elhub.devxp.build.configuration.pipeline.dsl.elhubProject
import no.elhub.devxp.build.configuration.pipeline.jobs.customJob
import no.elhub.devxp.build.configuration.pipeline.jobs.makeVerify
import no.elhub.devxp.build.configuration.pipeline.jobs.publishTag

elhubProject(DEVXP, "gh-dxp") {

    params {
        param("env.PATH", "\$PATH:/usr/local/go/bin:/usr/bin")
        param("env.GOROOT", "/usr/local/go")
    }

    pipeline {
        makeVerify {
            buildArtifactRules = listOf(ArtifactRule.include("build/coverage.*", "build.zip"))
            outputArtifactRules = listOf(ArtifactRule.include("build.zip!**", "build/"))
            sonarScanSettings = {
                sonarProjectSources = "."
                additionalParams = arrayListOf("-Dsonar.go.coverage.reportPaths=build/coverage.out")
            }
            enablePublishMetrics = true
        }
        publishTag()
        customJob(LinuxAgentContext) {
            name = "🚀 Release"
            id("Release")
            steps {
                script {
                    name = "GitHub Release Script"
                    scriptContent = """
                        make release
                    """.trimIndent()
                }
            }
        }
    }
}
