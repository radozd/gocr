SET root_dir=%cd%

mkdir build
cd build
mkdir bin

set INSTALL_DIR=%cd%

call vsvars64.bat

set OUT_INCLUDE_DIR=%INSTALL_DIR%\include\
set OUT_LIB_DIR=%INSTALL_DIR%\lib\

set INCLUDE=%INCLUDE%;%INSTALL_DIR%\include\
set LIBPATH=%LIBPATH%;%INSTALL_DIR%\lib\
set TESSDATA_PREFIX=%INSTALL_DIR%\share\tesseract\tessdata\

rem make STATIC
cd %INSTALL_DIR%
curl -O https://zlib.net/zlib131.zip
"c:\Program Files\Git\usr\bin\unzip.exe" zlib131.zip
cd zlib-1.3.1
cmake -Bbuild -DCMAKE_BUILD_TYPE=Release -DCMAKE_PREFIX_PATH=%INSTALL_DIR% -DCMAKE_INSTALL_PREFIX=%INSTALL_DIR% -DBUILD_SHARED_LIBS=False
cmake --build build --config Release --target zlibstatic
xcopy ".\Build\zconf.h" %OUT_INCLUDE_DIR% /y
xcopy ".\zlib.h"        %OUT_INCLUDE_DIR% /y
xcopy ".\zutil.h"       %OUT_INCLUDE_DIR% /y
xcopy ".\Build\Release\zlibstatic.lib" %OUT_LIB_DIR% /y
cd ..


rem make STATIC не умеет в BUILD_SHARED_LIB, но умеет в свои
cd %INSTALL_DIR%
curl -O -L https://downloads.sourceforge.net/project/libpng/libpng16/1.6.40/lpng1640.zip
"c:\Program Files\Git\usr\bin\unzip.exe" lpng1640.zip
cd lpng1640
cmake -Bbuild -DCMAKE_BUILD_TYPE=Release -DCMAKE_PREFIX_PATH=%INSTALL_DIR% -DCMAKE_INSTALL_PREFIX=%INSTALL_DIR% -DPNG_SHARED=False
cmake --build build --config Release --target install
cd ..


cd %INSTALL_DIR%
curl -O -L http://www.ijg.org/files/jpegsr9e.zip
"c:\Program Files\Git\usr\bin\unzip.exe" jpegsr9e.zip
cd jpeg-9e
nmake /f makefile.vs setup-v17
msbuild jpeg.sln
rem cmake --build build --config Release --target install
rem Так как не cmake, то просто копируем нужные нам файлы по папкам
xcopy ".\j*.h" "%INSTALL_DIR%\include\" /y
xcopy ".\Release\x64\jpeg.lib" "%INSTALL_DIR%\lib\jpeg.lib*" /y
cd ..


rem make STATIC умеет, но какая-то фигня
cd %INSTALL_DIR%
curl -O -L http://download.osgeo.org/libtiff/tiff-4.6.0.zip
"c:\Program Files\Git\usr\bin\unzip.exe" tiff-4.6.0.zip
cd tiff-4.6.0
cmake -Bbuild -DCMAKE_BUILD_TYPE=Release -DCMAKE_PREFIX_PATH=%INSTALL_DIR% -DCMAKE_INSTALL_PREFIX=%INSTALL_DIR% -DBUILD_SHARED_LIBS=False ^
    -Dtiff-tools=False ^
    -Dtiff-tests=False ^
    -Dtiff-contrib=False ^
    -Dtiff-docs=False
cmake --build build --config Release --target install
cd ..


rem make DLL
cd %INSTALL_DIR%
git clone --depth 1 https://github.com/DanBloomberg/leptonica.git
cd leptonica
cmake -Bbuild -DCMAKE_BUILD_TYPE=Release -DCMAKE_INSTALL_PREFIX=%INSTALL_DIR% -DCMAKE_PREFIX_PATH=%INSTALL_DIR% -DBUILD_PROG=OFF -DSW_BUILD=OFF -DBUILD_SHARED_LIBS=ON
cmake --build build  --config Release --target install
cd ..


rem make DLL
cd %INSTALL_DIR%
git clone --depth 1 https://github.com/tesseract-ocr/tesseract.git
cd tesseract
rem powershell "%root_dir%\add_spec_line.ps1" ".\CMakeLists.txt"
cmake -Bbuild -DCMAKE_BUILD_TYPE=Release -DCMAKE_INSTALL_PREFIX=%INSTALL_DIR%  -DCMAKE_PREFIX_PATH=%INSTALL_DIR% -DLeptonica_DIR=%INSTALL_DIR%\lib\cmake\leptonica -DBUILD_TRAINING_TOOLS=OFF -DSW_BUILD=OFF -DOPENMP_BUILD=OFF -DBUILD_SHARED_LIBS=ON
cmake --build build --config Release --target install
cd ..

cd ..