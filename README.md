# SMasteroids

[![Build Status](https://travis-ci.org/20zinnm/smasteroids.svg?branch=master)](https://travis-ci.org/20zinnm/smasteroids) [![codecov](https://codecov.io/gh/20zinnm/smasteroids/branch/master/graph/badge.svg)](https://codecov.io/gh/20zinnm/smasteroids)

> A game in honor of the St. Mark's science faculty and the Winn Science Building.

![Banner](docs/images/banner.png)

## Playing the Game

You can obtain pre-built releases of the game from the [Releases](https://github.com/20zinnm/smasteroids/releases) page. These are generally more reliable than development builds.

Note that for Windows, you must use 7-Zip to unzip the folder, as File Explorer does not work properly.

If your platform is not already included, you may also compile the program yourself using a standard Go toolchain with CGO enabled.

## Compatible Joysticks

As of right now, there are mappings for the following controllers:
* Joy-Con (L)
* Joy-Con (R)
* 8Bitdo SFC30 GamePad
* Dualshock 3 controller.

If you would like to use a different controller, either use a program to bind it to the keyboard controls, or open an issue with the specific type of controller.

## Built With

* [faiface/pixel](https://github.com/faiface/pixel) -- seriously, an incredible game library.
* [jakecoffman/cp](https://github.com/jakecoffman/cp) - 2D physics engine.
* Build scripts are based on MIT-licensed code by Humphrey Shotton.
* Blood, sweat, and tears.

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the tags on this repository.

## Authors

* Meyer Zinn - _project lead_ - Junior, St. Mark's

## License

This project is licensed under the GPLv3 License - see the LICENSE.md file for details

## Acknowledgements

* Faraz Asim - _playtester_ - Junior, St. Mark's
* Doug Rummel, who put up with this during Information Engineering.
* Fletcher Carron, who humored me by putting the game on the big screen in the Winn Science Center.

This game is meant to honor the science faculty in a humorous way. Quotes may be altered or fabricated for comedic purposes.
