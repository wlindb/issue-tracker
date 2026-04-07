<#import "template.ftl" as layout>

<@layout.registrationLayout displayMessage=!messagesPerField.existsError('username','password') displayInfo=realm.password && realm.registrationAllowed && !registrationDisabled??; section>

  <#if section == "header">
    <div class="cl-header">
      <div class="cl-logo-wrap">
        <svg class="cl-logo-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 40 40" fill="none" aria-hidden="true">
          <rect width="40" height="40" rx="10" fill="#111827"/>
          <path d="M11 8H29V12H22V28H29V32H11V28H18V12H11Z" fill="white"/>
        </svg>
      </div>
      <h1 class="cl-title">
        <#if realm.displayName?has_content>${realm.displayName}<#else>Sign in</#if>
      </h1>
      <p class="cl-subtitle">
        <#if realm.displayNameHtml?has_content>
          ${realm.displayNameHtml?no_esc}
        <#else>
          Welcome back! Please sign in to continue.
        </#if>
      </p>
    </div>
  </#if>

  <#if section == "form">
    <div id="kc-form" class="cl-form-wrap">

      <#-- Social / Identity Provider login -->
      <#if realm.password && social.providers??>
        <div id="kc-social-providers" class="cl-social-providers">
          <#list social.providers as p>
            <a id="social-${p.alias}" class="cl-social-btn" href="${p.loginUrl}" type="button">
              <#if p.iconClasses?has_content>
                <i class="${p.iconClasses!}" aria-hidden="true"></i>
              </#if>
              <span>Continue with ${p.displayName!}</span>
            </a>
          </#list>
        </div>

        <#if realm.password>
          <div class="cl-divider">
            <span class="cl-divider-line"></span>
            <span class="cl-divider-label">or</span>
            <span class="cl-divider-line"></span>
          </div>
        </#if>
      </#if>

      <#-- Username / Password Form -->
      <#if realm.password>
        <form id="kc-form-login" class="cl-form" onsubmit="document.getElementById('kc-login').disabled = true; return true;" action="${url.loginAction}" method="post">
          <div class="cl-field-group">
            <#-- Username / Email -->
            <div class="cl-field ${messagesPerField.printIfExists('username', 'cl-field--error')}">
              <label class="cl-label" for="username">
                <#if !realm.loginWithEmailAllowed>${msg("username")}
                <#elseif !realm.registrationEmailAsUsername>${msg("usernameOrEmail")}
                <#else>${msg("email")}</#if>
              </label>
              <input
                tabindex="1"
                id="username"
                class="cl-input"
                name="username"
                type="text"
                autofocus
                autocomplete="username"
                value="${(login.username!'')}"
                aria-invalid="<#if messagesPerField.existsError('username','password')>true</#if>"
              />
              <#if messagesPerField.existsError('username','password')>
                <span class="cl-field-error" aria-live="polite">
                  ${kcSanitize(messagesPerField.getFirstError('username','password'))?no_esc}
                </span>
              </#if>
            </div>

            <#-- Password -->
            <div class="cl-field ${messagesPerField.printIfExists('password', 'cl-field--error')}">
              <div class="cl-label-row">
                <label class="cl-label" for="password">${msg("password")}</label>
                <#if realm.resetPasswordAllowed>
                  <a tabindex="5" class="cl-forgot-link" href="${url.loginResetCredentialsUrl}">${msg("doForgotPassword")}</a>
                </#if>
              </div>
              <div class="cl-input-wrap">
                <input
                  tabindex="2"
                  id="password"
                  class="cl-input"
                  name="password"
                  type="password"
                  autocomplete="current-password"
                  aria-invalid="<#if messagesPerField.existsError('username','password')>true</#if>"
                />
                <button type="button" class="cl-password-toggle" aria-label="Toggle password visibility" onclick="togglePassword()">
                  <svg id="eye-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
                    <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
                    <circle cx="12" cy="12" r="3"/>
                  </svg>
                  <svg id="eye-off-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true" style="display:none">
                    <path d="M17.94 17.94A10.07 10.07 0 0112 20c-7 0-11-8-11-8a18.45 18.45 0 015.06-5.94M9.9 4.24A9.12 9.12 0 0112 4c7 0 11 8 11 8a18.5 18.5 0 01-2.16 3.19m-6.72-1.07a3 3 0 11-4.24-4.24"/>
                    <line x1="1" y1="1" x2="23" y2="23"/>
                  </svg>
                </button>
              </div>
            </div>
          </div>

          <#-- Remember Me -->
          <#if realm.rememberMe && !usernameEditDisabled??>
            <div class="cl-remember-wrap">
              <label class="cl-checkbox-label">
                <input
                  tabindex="3"
                  id="rememberMe"
                  name="rememberMe"
                  type="checkbox"
                  class="cl-checkbox"
                  <#if login.rememberMe??>checked</#if>
                />
                <span class="cl-checkbox-custom"></span>
                <span>${msg("rememberMe")}</span>
              </label>
            </div>
          </#if>

          <input type="hidden" id="id-hidden-input" name="credentialId" <#if auth.selectedCredential?has_content>value="${auth.selectedCredential}"</#if>/>

          <button tabindex="4" class="cl-primary-btn" name="login" id="kc-login" type="submit">
            ${msg("doLogIn")}
          </button>
        </form>
      </#if>
    </div>
  </#if>

  <#if section == "info">
    <#if realm.password && realm.registrationAllowed && !registrationDisabled??>
      <div class="cl-footer">
        <span>${msg("noAccount")}</span>
        <a tabindex="6" class="cl-footer-link" href="${url.registrationUrl}">${msg("doRegister")}</a>
      </div>
    </#if>
  </#if>

</@layout.registrationLayout>

<script>
  function togglePassword() {
    const input = document.getElementById('password');
    const eyeIcon = document.getElementById('eye-icon');
    const eyeOffIcon = document.getElementById('eye-off-icon');
    if (input.type === 'password') {
      input.type = 'text';
      eyeIcon.style.display = 'none';
      eyeOffIcon.style.display = 'block';
    } else {
      input.type = 'password';
      eyeIcon.style.display = 'block';
      eyeOffIcon.style.display = 'none';
    }
  }
</script>
