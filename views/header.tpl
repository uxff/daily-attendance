<header id="topbar" class="navbar navbar-default navbar-fixed-top bs-docs-nav" role="banner">
  <div class="container">
    <div class="row">
    <div class="navbar-header">
      <button class="navbar-toggle collapsed" type="button" data-toggle="collapse" data-target=".bs-navbar-collapse">
        <span class="sr-only">导航</span>
        <span class="icon-bar"></span>
        <span class="icon-bar"></span>
        <span class="icon-bar"></span>
      </button>
      <a style="font-size: 14px;" class="navbar-brand" rel="home" href="/" >
        <strong>{{.appname}}</strong>
      </a>
    </div>

    <nav class="collapse navbar-collapse bs-navbar-collapse" role="navigation" >
      <ul class="nav navbar-nav">
        <li><a href='{{urlfor "IndexController.Index"}}'>
          <span class="glyphicon glyphicon-home"></span> 首页
        </a></li>
        <li><a href='{{urlfor "IndexController.Index"}}'>
          <span class="glyphicon glyphicon-heart"></span> 关注健康
        </a></li>
        <li>
          <a href="javascript:;" class="dropdown-toggle" data-hover="dropdown">
          <span class="glyphicon glyphicon-time"></span> 定制计划 <b class="caret"></b>
          </a>
            <ul class="dropdown-menu">
                <li role="presentation" class="dropdown-header">计划介绍</li>
                <li><a href="/attendance">每日打卡</a></li>
                <li><a href="/">瓜分气馁者</a></li>
                <li role="presentation" class="divider"></li>
                <li role="presentation" class="dropdown-header">排行</li>
                <li><a href="/">坚持最长排行</a></li>
                <li><a href="/">获得瓜分最多</a></li>
            </ul>
        </li>
        <li><a href='{{urlfor "IndexController.Index"}}'>
          <span class="glyphicon glyphicon-list-alt"></span> 我的计划
        </a></li>
      </ul>

      <ul class="nav navbar-nav navbar-right">
        <li class="dropdown">
          <a href="{{urlfor "UsersController.Index"}}" role="button" class="dropdown-toggle" data-hover="dropdown">
            <span class='glyphicon glyphicon-user'></span> 用户{{if .IsLogin}}({{.Userinfo.Email}}){{end}} <b class="caret"></b>
          </a>
          <ul class="dropdown-menu">
            {{if .IsLogin}}
                <li ><a href='{{urlfor "UsersController.Logout"}}'>
                  <span class='glyphicon glyphicon-log-out'></span> 退出
                </a></li>
            {{else}}
                <li ><a href='{{urlfor "UsersController.Login"}}'>
                  <span class='glyphicon glyphicon-log-in'></span> 登录
                </a></li>
                <li ><a href='{{urlfor "UsersController.Signup"}}'>
                    <span class='glyphicon glyphicon-check'></span> 注册
                </a></li>
            {{end}}
          </ul>
        </li>
      </ul>
    </nav>
    </div>
  </div>

</header>
