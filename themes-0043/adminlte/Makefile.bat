SET CLI=adm
SET RESOURCE_PATH=.\resource
SET ASSETS_PATH=.\resource\assets
SET SEPARATION_PATH=.\separation
SET COMMON_PATH=.\..\common

DEL /Q /F /S %ASSETS_PATH%\dist
mkdir %ASSETS_PATH%\dist
mkdir %ASSETS_PATH%\dist\js
mkdir %ASSETS_PATH%\dist\css
xcopy /Y /E %COMMON_PATH%\assets\* %ASSETS_PATH%\src\
copy /Y %RESOURCE_PATH%\adminlte\adminlte.css %ASSETS_PATH%\src\css\combine\g_all.css
xcopy /Y /E %ASSETS_PATH%\src\js\*.js %ASSETS_PATH%\dist\js\
xcopy /Y /E %ASSETS_PATH%\src\css\*.png %ASSETS_PATH%\dist\css\
xcopy /Y /E %COMMON_PATH%\pages\* %RESOURCE_PATH%\pages\

xcopy /Y /E %ASSETS_PATH%\src\css\fonts %ASSETS_PATH%\dist\css\fonts\
xcopy /Y /E %ASSETS_PATH%\src\img %ASSETS_PATH%\dist\img\
xcopy /Y /E %ASSETS_PATH%\src\fonts %ASSETS_PATH%\dist\fonts\

DEL /Q /F /S %SEPARATION_PATH%\public
mkdir %SEPARATION_PATH%\public

xcopy /Y /E %RESOURCE_PATH%\assets %SEPARATION_PATH%\public\assets
xcopy /Y /E %RESOURCE_PATH%\pages %SEPARATION_PATH%\public\pages
DEL /Q /F /S %SEPARATION_PATH%\public\assets\vendor


%CLI% combine js --hash=true --src=%ASSETS_PATH%\src\js\all\ --dist=%ASSETS_PATH%\dist\js\all.min.js
%CLI% combine js --hash=true --src=%ASSETS_PATH%\src\js\all_2\ --dist=%ASSETS_PATH%\dist\js\all_2.min.js
%CLI% combine js --hash=true --src=%ASSETS_PATH%\src\js\components\form\ --dist=%ASSETS_PATH%\dist\js\form.min.js
%CLI% combine js --hash=true --src=%ASSETS_PATH%\src\js\components\tree\ --dist=%ASSETS_PATH%\dist\js\tree.min.js
%CLI% combine js --hash=true --src=%ASSETS_PATH%\src\js\components\treeview\ --dist=%ASSETS_PATH%\dist\js\treeview.min.js
%CLI% combine js --hash=true --src=%ASSETS_PATH%\src\js\components\datatable\ --dist=%ASSETS_PATH%\dist\js\datatable.min.js
xcopy /Y /E %ASSETS_PATH%\dist\js\* %SEPARATION_PATH%\public\assets\dist\js\

%CLI% combine css --hash=true
xcopy /Y /E %ASSETS_PATH%\dist\css\*.css %SEPARATION_PATH%\public\assets\dist\css\

%CLI% compile asset
packr2 clean
packr2

%CLI% compile tpl -p=adminlte
    
DEL /Q /F /S %ASSETS_PATH%\src\*
DEL /Q /F /S %RESOURCE_PATH%\pages\*

copy /Y adminlte-packr.txt adminlte-packr.go

go fmt .\...
pause
