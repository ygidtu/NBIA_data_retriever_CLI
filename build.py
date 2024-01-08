#!/usr/bin/env python
import os
from argparse import ArgumentParser
from datetime import datetime
from subprocess import check_output, check_call, CalledProcessError


__VERSION__ = '1.1.0'


def __init_name__():
    GOMOD = os.path.join(os.path.dirname(__file__), "go.mod")

    if not os.path.exists(GOMOD):
        raise FileNotFoundError("go.mod file not found under same directory of this script")

    with open(GOMOD) as r:
        for line in r:
            if line.startswith("module"):
                line = line.split()
                GOMOD = os.path.basename(line[-1])
                break

    return GOMOD if not GOMOD.endswith("go.mod") else "build"


def __init_architectures__():
    platforms = {}
    architectures = set()
    for x in check_output("go tool dist list", shell=True).decode("utf-8").split("\n"):
        if not x:
            continue
        k, v = os.path.dirname(x), os.path.basename(x)

        if k not in platforms.keys():
            platforms[k] = []
        platforms[k].append(v)
        architectures.add(v)
    return platforms, architectures


__ARCHITECTURE__, ARCHITECTURE = __init_architectures__()
GOMOD = __init_name__()
PLATFORMS = sorted(__ARCHITECTURE__.keys())


def main():

    parser = ArgumentParser(f"build {GOMOD} from source code")
    parser.add_argument("-p", "--platform", default=None,
                        help="which platform to build, one of {}".format(",".join(PLATFORMS)), type=str)
    parser.add_argument("-a", "--arch", default=None,
                        help="which architecture to build, one of {}".format(",".join(ARCHITECTURE)), type=str)
    parser.add_argument("--all", action='store_true',
                        help="build for all supported architectures")
    parser.add_argument("--common", action='store_true',
                        help="build for windows, linux and darwin architectures")

    args = parser.parse_args()

    if not args.platform and not args.arch and not args.all and not args.common:
        print("info: -h/--help print for usage information")

    platform = [check_output("go env GOHOSTOS", shell=True).decode("utf-8").strip()]
    arch = [check_output("go env GOHOSTARCH", shell=True).decode("utf-8").strip()]

    if args.all:
        print("build for all architectures")
        platform = PLATFORMS
        arch = ARCHITECTURE
    elif not args.all:
        if args.common:
            print("build for most common architectures")
            platform = [x for x in PLATFORMS if x in ["windows", "linux", "darwin"]]
            arch = [x for x in ARCHITECTURE if x in ["amd64", "arm64"]]
        elif args.platform or args.arch:
            if args.platform:
                if args.platform.lower() in PLATFORMS:
                    platform = [args.platform.lower()]
                else:
                    raise ValueError("{} is not supported".format(args.platform))

            if args.arch:
                if args.arch.lower() in ARCHITECTURE:
                    arch = [args.arch.lower()]
                else:
                    raise ValueError("{} is not supported".format(args.arch))
        else:
            print("build for current platform")

    date = datetime.now().strftime("%Y-%m-%d")
    gitVersion = check_output("git rev-parse HEAD", shell=True).decode("utf-8").strip()
    goVersion = ""
    for i in check_output("go version", shell=True).decode("utf-8").strip().split():
        if i.strip() != "go" and i.startswith("go"):
            goVersion = i.strip().replace("go", "")

    flags = f"-X main.buildStamp={date} -X main.gitHash={gitVersion} -X main.goVersion={goVersion} -X main.version={__VERSION__} -s -w"

    for i in platform:
        for j in arch:
            if j in __ARCHITECTURE__[i]:
                try:
                    print("building for: {} {}".format(i, j))

                    with open(os.devnull, "w") as w:
                        check_call(f"env GOOS='{i}' GOARCH={j} go build -ldflags \"{flags}\" -x -o {GOMOD}_{i}_{j} .", shell=True, stdout=w, stderr=w)
                except CalledProcessError:
                    continue


if __name__ == "__main__":
    main()
