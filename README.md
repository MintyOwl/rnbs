# rnbs
Bootstrap your boilerplate RN project

## Install from binary
Binaries will be available (soon)

## Install from source
* Requires go >= 1.8

#### Once go is installed
    go get github.com/MintyOwl/rnbs

    go install (will install in PATH)

#### Now assuming you already have performed *react-native init MyAwesomeProject* and edited it to add boilerplate

    rnbs.exe path/to/MyAwesomeProject path/to/BrandNewAmazingProject

* Where ***BrandNewAmazingProject*** is a folder/direcotory that doesnot exist, will be created for you with all your existing boilpate from ***MyAwesomeProject***

* Now just do *react-native run-android* inside ***BrandNewAmazingProject***

Note:

* I dont have MacOS, so ios is not tested. However, theoritically you should be able to run *react-native run-ios*

* I need to know what build folders are created for ios. For ex: in android (on Windows 10), you have .gradle, androis/app/build and android/build folders after run-android.