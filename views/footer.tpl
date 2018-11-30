<footer>
  <div class="container">
    <div class="clearfix">
      <div class="footer-logo">
        <a href="/">
          <img src=""><small style='font-size: 65%;'>{{.appname}}</small>
        </a>
      </div>
      <dl class="footer-nav">
        <dt class="nav-title">关注健康</dt>
        <dd class="nav-item">
          <a href="#">
            <span class="glyphicon glyphicon-credit-card"> 健康资讯 </span>
          </a>
        </dd>
        <dd class="nav-item">
          <a href="#" target="_blank">
            <span class='glyphicon glyphicon-bullhorn'> 定制计划 </span>
          </a>
        </dd>
      </dl>
      <dl class="footer-nav">
        <dt class="nav-title">关于</dt>

        <dd class="nav-item">
          <a href="#">
            <span class='glyphicon glyphicon-info-sign'> 关于我们 </span>
          </a>
        </dd>

      </dl>

      <dl class="footer-nav hidden">
        <dt class="nav-title">社会反馈</dt>
        <dd class="nav-item">
          <a href="#" target="_blank">
            <span class='glyphicon glyphicon-globe'></span> 
          </a>
        </dd>
      </dl>

      <dl class="footer-nav">
        <dt class="nav-title">联系方式</dt>
        <dd class="nav-item">
          <a href="#">
            <span class='glyphicon glyphicon-comment'> 联系站长 </span>
          </a>
        </dd>
        <dd class="nav-item">
          <a href="#">
            <span class='glyphicon glyphicon-comment'> 微信 </span>
          </a>
        </dd>
        <dd class="nav-item">
          <a href="#">
            <span class='glyphicon glyphicon-comment'> QQ群: 12345678 </span>
          </a>
        </dd>
      </dl>

    </div>

    <div class="footer-copyright text-center">
      友情链接:
      {{range $k, $link := .friendlyLinks}}
          <a href="{{$link.Url}}" target="_blank">{{$link.Name}}</a> &nbsp;
      {{end}}
    </div>
    <div class="footer-copyright text-center">
      Copyright <span class="glyphicon glyphicon-copyright-mark"></span>
      2014-{{datenow "2006"}} <strong>{{.appname}}</strong>
      All rights reserved.
    </div>

  </div>
</footer>
