package lilypond

import (
	"fmt"
	"github.com/lsierant/notes-gen/pkg/notes"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type Renderer struct {
	WorkingDir string
}

func (r *Renderer) RenderPNG(source string) ([]byte, error) {
	tmpDir, err := ioutil.TempDir(r.WorkingDir, "sources-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %v", tmpDir)
	}

	outDir := fmt.Sprintf("%s/out", tmpDir)
	err = os.Mkdir(outDir, 0770)
	if err != nil {
		return nil, fmt.Errorf("failed to create out dir %s: %v", outDir, err)
	}

	filename := "1.ly"
	err = ioutil.WriteFile(fmt.Sprintf("%s/%s", tmpDir, filename), []byte(source), 0660)
	if err != nil {
		return nil, fmt.Errorf("failed to write source file: %s/%s: %v", tmpDir, filename, err)
	}

	pwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working dir: %v", err)
	}

	argsStr := `run -v %s/%s:/d docker.io/airdock/lilypond:latest -dresolution=300 --png -dbackend=eps -dno-gs-load-fonts -dinclude-eps-fonts -o /d/out /d/1.ly`
	args := strings.Split(fmt.Sprintf(argsStr, pwd, tmpDir), " ")
	commandName := "docker"
	command := exec.Command(commandName, args...)
	output, err := command.CombinedOutput()
	fmt.Printf("running command: \n%s %s\n", commandName, strings.Join(args, " "))
	fmt.Printf("%s", output)
	if err != nil {
		return nil, fmt.Errorf("error running command %s %s: %v", commandName, strings.Join(args, " "), err)
	}

	pngBytes, err := ioutil.ReadFile(fmt.Sprintf("%s/1.png", outDir))
	if err != nil {
		return nil, fmt.Errorf("failed to read png file: %v", err)
	}

	err = os.RemoveAll(tmpDir)
	if err != nil {
		return nil, fmt.Errorf("failed to cleanup tmp dir %s: %v", tmpDir, err)
	}
	return pngBytes, nil
}

func RenderIntervalImage(renderer *Renderer, interval notes.Interval) ([]byte, error) {
	first := interval.FirstNote
	second := interval.SecondNote

	var source string
	var err error
	if first.BassClef && second.BassClef {
		source, err = parseAndRenderTextTemplate("interval", bassOnlyIntervalTemplate, interval)
	} else if first.BassClef && second.TrebleClef {
		source, err = parseAndRenderTextTemplate("interval", bassAndTrebleIntervalTemplate, interval)
	} else if first.TrebleClef && second.TrebleClef {
		source, err = parseAndRenderTextTemplate("interval", trebleOnlyIntervalTemplate, interval)
	} else {
		return nil, fmt.Errorf("not supported interval: %+v", interval)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to render interval template: %v", err)
	}

	png, err := renderer.RenderPNG(source)
	if err != nil {
		return nil, fmt.Errorf("failed to render PNG from source: %s: %v", source, err)
	}

	return png, err
}
