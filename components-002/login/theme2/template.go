package theme2

var List = map[string]string{"login": `{{define "login"}}
    <!DOCTYPE html>
    <!--[if lt IE 7]>
    <html class="no-js lt-ie9 lt-ie8 lt-ie7">
    <![endif]-->
    <!--[if IE 7]>
    <html class="no-js lt-ie9 lt-ie8">
    <![endif]-->
    <!--[if IE 8]>
    <html class="no-js lt-ie9">
    <![endif]-->
    <!--[if gt IE 8]><!-->
    <html class="no-js">
    <!--<![endif]-->
    <head>
        <meta charset="utf-8">
        <meta http-equiv="X-UA-Compatible" content="IE=edge">
        <title>{{.Title}}</title>
        <meta name="viewport" content="width=device-width, initial-scale=1">

        <link rel="stylesheet" href="{{link .CdnUrl .UrlPrefix "/assets/login/dist/all.min.css"}}">

        <!--[if lt IE 9]>
        <script src="{{link .CdnUrl .UrlPrefix "/assets/login/dist/respond.min.js"}}"></script>
        <![endif]-->
    </head>
    <body>

    <div class="container">
        <div class="row" style="margin-top: 80px;">
            <div class="col-md-4 col-md-offset-4">
                <form onsubmit="return false" action="##" method="post" id="sign-up-form" class="fh5co-form animate-box"
                      data-animate-effect="fadeIn">
                    <h2>{{.Title}}</h2>
                    <div class="form-group">
                        <label for="username" class="sr-only">{{lang "username"}}</label>
                        <input type="text" required class="form-control" id="username" placeholder="{{lang "username"}}"
                               autocomplete="off" value="">
                    </div>
                    <div class="form-group">
                        <label for="password" class="sr-only">{{lang "password"}}</label>
                        <input type="password" required class="form-control" id="password" placeholder="{{lang "password"}}"
                               autocomplete="off" value="">
                    </div>
                    {% if .CaptchaDigits %}
                    <div class="form-group has-feedback 1">
                        <div class="row">
                            <div class="col-xs-7">
                                <input type="text" class="form-control" placeholder="{{lang "captcha"}}" id="captcha">
                            </div>
                            <div class="col-xs-5">
                                <img id="captcha-img" class="captcha" src="{% .CaptchaImgSrc %}" alt="" width="110" height="45">
                            </div>
                        </div>
                        <input type="hidden" value="{% .CaptchaID %}" id="captcha-id">
                    </div>
                    {% end %}
                    <div class="form-group">
                        <button class="btn btn-primary" onclick="submitData()">{{lang "login"}}</button>
                    </div>
                </form>
            </div>
        </div>
        <div class="row" style="padding-top: 60px; clear: both;">
            <div class="col-md-12 text-center">
            </div>
        </div>
    </div>

    <div id="particles-js">
        <canvas class="particles-js-canvas-el" width="1606" height="1862" style="width: 100%; height: 100%;"></canvas>
    </div>

    <script src="{{link .CdnUrl .UrlPrefix "/assets/login/dist/all.min.js"}}"></script>

    {% if .TencentWaterProofWallData.AppID  %}
        <script src="https://ssl.captcha.qq.com/TCaptcha.js"></script>
    {% end %}

    <script>

        {% if .TencentWaterProofWallData.AppID  %}

        let captcha = new TencentCaptcha("{% .TencentWaterProofWallData.AppID %}", function (res) {
            // res（用户主动关闭验证码）= {ret: 2, ticket: null}
            // res（验证成功） = {ret: 0, ticket: "String", randstr: "String"}
            if (res.ret === 0) {
                $.ajax({
                    dataType: 'json',
                    type: 'POST',
                    url: '{{.UrlPrefix}}/signin',
                    async: 'true',
                    data: {
                        'username': $("#username").val(),
                        'password': $("#password").val(),
                        'token': res.ticket+","+res.randstr
                    },
                    success: function (data) {
                        location.href = data.data.url;
                    },
                    error: function (data) {
                        alert(data.responseJSON.msg);
                    }
                });
            } else {
                alert(data.data.msg);
            }
        }, {});

        {% end %}

        {% if .CaptchaDigits %}

        $("#captcha-img").on("click",function(){
            // location.reload();
            $.ajax({
                    dataType: 'json',
                    type: 'POST',
                    url: '/captchaRefresh',
                    async: 'true',
                    data: {},
                    success: function (data) {
                        $("#captcha-id").val(data.id);
                        $("#captcha-img").attr('src',data.img);
                    },
                    error: function (data) {
                        alert(data.responseJSON.msg);
                    }
                });
        }).click();

        {% end %}

        function submitData() {
            {% if .TencentWaterProofWallData.AppID  %}
            captcha.show();
            {% else  %}
            $.ajax({
                dataType: 'json',
                type: 'POST',
                url: '{{.UrlPrefix}}/signin',
                async: 'true',
                data: {
                    'username': $("#username").val(),
                    {% if .CaptchaDigits %}
                    'token': $("#captcha").val() + "," + $("#captcha-id").val(),
                    {% end %}
                    'password': $("#password").val()
                },
                success: function (data) {
                    location.href = data.data.url
                },
                error: function (data) {
                    alert(data.responseJSON.msg);
                    $("#captcha-img").click();
                }
             });
            {% end %}
        }
    </script>

    </body>
    </html>
{{end}}`}
