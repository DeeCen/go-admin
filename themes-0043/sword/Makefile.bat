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
    xcopy /Y /E %RESOURCE_PATH%\sword\sword.css %ASSETS_PATH%\src\css\combine\i_sword.css
    xcopy /Y %RESOURCE_PATH%\sword\blue@2x.png %ASSETS_PATH%\src\css\
    xcopy /Y %RESOURCE_PATH%\sword\blue.png %ASSETS_PATH%\src\css\

    xcopy /Y %ASSETS_PATH%\src\js\*.js %ASSETS_PATH%\dist\js\
    xcopy /Y %ASSETS_PATH%\src\css\*.png %ASSETS_PATH%\dist\css\

    xcopy /Y /E %COMMON_PATH%\pages\* %RESOURCE_PATH%\pages\
    xcopy /Y %RESOURCE_PATH%\sword\pages\*.tmpl %RESOURCE_PATH%\pages\
    xcopy /Y %RESOURCE_PATH%\sword\pages\components\*.tmpl %RESOURCE_PATH%\pages\components\
    xcopy /Y %RESOURCE_PATH%\sword\pages\components\table\*.tmpl %RESOURCE_PATH%\pages\components\table\

    xcopy /Y /E %ASSETS_PATH%\src\css\fonts %ASSETS_PATH%\dist\css\
    xcopy /Y /E %ASSETS_PATH%\src\img %ASSETS_PATH%\dist\
    xcopy /Y /E %ASSETS_PATH%\src\fonts %ASSETS_PATH%\dist\

    xcopy /Y /E %RESOURCE_PATH%\assets %SEPARATION_PATH%\public\assets
    xcopy /Y /E %RESOURCE_PATH%\pages %SEPARATION_PATH%\public\pages
    DEL /Q /F /S %SEPARATION_PATH%\public\assets\vendor
	
	
	
	
	%CLI% combine js --hash=true --src=%ASSETS_PATH%\src\js\all\ --dist=%ASSETS_PATH%\dist\js\all.min.js
    %CLI% combine js --hash=true --src=%ASSETS_PATH%\src\js\all_2\ --dist=%ASSETS_PATH%\dist\js\all_2.min.js
    %CLI% combine js --hash=true --src=%ASSETS_PATH%\src\js\components\form\ --dist=%ASSETS_PATH%\dist\js\form.min.js
    %CLI% combine js --hash=true --src=%ASSETS_PATH%\src\js\components\tree\ --dist=%ASSETS_PATH%\dist\js\tree.min.js
    %CLI% combine js --hash=true --src=%ASSETS_PATH%\src\js\components\treeview\ --dist=%ASSETS_PATH%\dist\js\treeview.min.js
    %CLI% combine js --hash=true --src=%ASSETS_PATH%\src\js\components\datatable\ --dist=%ASSETS_PATH%\dist\js\datatable.min.js
    xcopy /Y %ASSETS_PATH%\dist\js\* %SEPARATION_PATH%\public\assets\dist\js\


    %CLI% combine css --hash=true
    xcopy /Y %ASSETS_PATH%\dist\css\*.css %SEPARATION_PATH%\public\assets\dist\css\
	
	
	%CLI% compile asset
    packr2 clean
    packr2
	
	
	%CLI% compile tpl -p=sword
	
	
    DEL /Q /F /S %ASSETS_PATH%\dist
    mkdir %ASSETS_PATH%\dist
    mkdir %ASSETS_PATH%\dist\js
    mkdir %ASSETS_PATH%\dist\css

    xcopy /Y /E %COMMON_PATH%\assets\* %ASSETS_PATH%\src\
    xcopy /Y %RESOURCE_PATH%\sword\sword.css %ASSETS_PATH%\src\css\combine\i_sword.css
    xcopy /Y %RESOURCE_PATH%\sword\blue@2x.png %ASSETS_PATH%\src\css\
    xcopy /Y %RESOURCE_PATH%\sword\blue.png %ASSETS_PATH%\src\css\

    xcopy /Y %ASSETS_PATH%\src\js\*.js %ASSETS_PATH%\dist\js\
    xcopy /Y %ASSETS_PATH%\src\css\*.png %ASSETS_PATH%\dist\css\

    xcopy /Y /E %COMMON_PATH%\pages\* %RESOURCE_PATH%\pages\
    xcopy /Y %RESOURCE_PATH%\sword\pages\*.tmpl %RESOURCE_PATH%\pages\
    xcopy /Y %RESOURCE_PATH%\sword\pages\components\*.tmpl %RESOURCE_PATH%\pages\components\
    xcopy /Y %RESOURCE_PATH%\sword\pages\components\table\*.tmpl %RESOURCE_PATH%\pages\components\table\

    xcopy /Y /E %ASSETS_PATH%\src\css\fonts %ASSETS_PATH%\dist\css\
    xcopy /Y /E %ASSETS_PATH%\src\img %ASSETS_PATH%\dist\
    xcopy /Y /E %ASSETS_PATH%\src\fonts %ASSETS_PATH%\dist\

    DEL /Q /F /S %SEPARATION_PATH%\public
    mkdir %SEPARATION_PATH%\public

    xcopy /Y /E %RESOURCE_PATH%\assets %SEPARATION_PATH%\public\assets
    xcopy /Y /E %RESOURCE_PATH%\pages %SEPARATION_PATH%\public\pages
    DEL /Q /F /S %SEPARATION_PATH%\public\assets\vendor

	DEL /Q /F /S %ASSETS_PATH%\src\*
    DEL /Q /F /S %RESOURCE_PATH%\pages\*
    
	go fmt .\...
	
	pause