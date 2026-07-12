<#import "template.ftl" as layout>

<@layout.registrationLayout displayMessage=!messagesPerField.existsError('firstName','lastName','email','username','password','password-confirm'); section>

  <#if section == "header">
    <div class="cl-header">
      <div class="cl-logo-wrap">
        <svg class="cl-logo-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 40 40" fill="none" aria-hidden="true">
          <rect width="40" height="40" rx="10" fill="#111827"/>
          <path d="M11 8H29V12H22V28H29V32H11V28H18V12H11Z" fill="white"/>
        </svg>
      </div>
      <h1 class="cl-title">Create an account</h1>
      <p class="cl-subtitle">Fill in the details below to get started.</p>
    </div>
  </#if>

  <#if section == "form">
    <div id="kc-register-form" class="cl-form-wrap">
      <form id="kc-register" class="cl-form" action="${url.registrationAction}" method="post">

        <div class="cl-field-group">

          <#-- First Name -->
          <div class="cl-field ${messagesPerField.printIfExists('firstName', 'cl-field--error')}">
            <label class="cl-label" for="firstName">${msg("firstName")}</label>
            <input
              tabindex="1"
              id="firstName"
              class="cl-input"
              name="firstName"
              type="text"
              autofocus
              autocomplete="given-name"
              value="${(register.formData.firstName!'')}"
              aria-invalid="<#if messagesPerField.existsError('firstName')>true</#if>"
            />
            <#if messagesPerField.existsError('firstName')>
              <span class="cl-field-error" aria-live="polite">
                ${kcSanitize(messagesPerField.get('firstName'))?no_esc}
              </span>
            </#if>
          </div>

          <#-- Last Name -->
          <div class="cl-field ${messagesPerField.printIfExists('lastName', 'cl-field--error')}">
            <label class="cl-label" for="lastName">${msg("lastName")}</label>
            <input
              tabindex="2"
              id="lastName"
              class="cl-input"
              name="lastName"
              type="text"
              autocomplete="family-name"
              value="${(register.formData.lastName!'')}"
              aria-invalid="<#if messagesPerField.existsError('lastName')>true</#if>"
            />
            <#if messagesPerField.existsError('lastName')>
              <span class="cl-field-error" aria-live="polite">
                ${kcSanitize(messagesPerField.get('lastName'))?no_esc}
              </span>
            </#if>
          </div>

          <#-- Email -->
          <div class="cl-field ${messagesPerField.printIfExists('email', 'cl-field--error')}">
            <label class="cl-label" for="email">${msg("email")}</label>
            <input
              tabindex="3"
              id="email"
              class="cl-input"
              name="email"
              type="email"
              autocomplete="email"
              value="${(register.formData.email!'')}"
              aria-invalid="<#if messagesPerField.existsError('email')>true</#if>"
            />
            <#if messagesPerField.existsError('email')>
              <span class="cl-field-error" aria-live="polite">
                ${kcSanitize(messagesPerField.get('email'))?no_esc}
              </span>
            </#if>
          </div>

          <#-- Username (only shown when email is not used as username) -->
          <#if !realm.registrationEmailAsUsername>
            <div class="cl-field ${messagesPerField.printIfExists('username', 'cl-field--error')}">
              <label class="cl-label" for="username">${msg("username")}</label>
              <input
                tabindex="4"
                id="username"
                class="cl-input"
                name="username"
                type="text"
                autocomplete="username"
                value="${(register.formData.username!'')}"
                aria-invalid="<#if messagesPerField.existsError('username')>true</#if>"
              />
              <#if messagesPerField.existsError('username')>
                <span class="cl-field-error" aria-live="polite">
                  ${kcSanitize(messagesPerField.get('username'))?no_esc}
                </span>
              </#if>
            </div>
          </#if>

          <#-- Password -->
          <#if passwordRequired??>
            <div class="cl-field ${messagesPerField.printIfExists('password', 'cl-field--error')}">
              <label class="cl-label" for="password">${msg("password")}</label>
              <div class="cl-input-wrap">
                <input
                  tabindex="5"
                  id="password"
                  class="cl-input"
                  name="password"
                  type="password"
                  autocomplete="new-password"
                  aria-invalid="<#if messagesPerField.existsError('password','password-confirm')>true</#if>"
                />
                <button type="button" class="cl-password-toggle" aria-label="Toggle password visibility" onclick="togglePassword('password','eye-icon','eye-off-icon')">
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
              <#if messagesPerField.existsError('password')>
                <span class="cl-field-error" aria-live="polite">
                  ${kcSanitize(messagesPerField.get('password'))?no_esc}
                </span>
              </#if>
            </div>

            <#-- Password Confirm -->
            <div class="cl-field ${messagesPerField.printIfExists('password-confirm', 'cl-field--error')}">
              <label class="cl-label" for="password-confirm">${msg("passwordConfirm")}</label>
              <div class="cl-input-wrap">
                <input
                  tabindex="6"
                  id="password-confirm"
                  class="cl-input"
                  name="password-confirm"
                  type="password"
                  autocomplete="new-password"
                  aria-invalid="<#if messagesPerField.existsError('password-confirm')>true</#if>"
                />
                <button type="button" class="cl-password-toggle" aria-label="Toggle confirm password visibility" onclick="togglePassword('password-confirm','eye-icon-confirm','eye-off-icon-confirm')">
                  <svg id="eye-icon-confirm" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
                    <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
                    <circle cx="12" cy="12" r="3"/>
                  </svg>
                  <svg id="eye-off-icon-confirm" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true" style="display:none">
                    <path d="M17.94 17.94A10.07 10.07 0 0112 20c-7 0-11-8-11-8a18.45 18.45 0 015.06-5.94M9.9 4.24A9.12 9.12 0 0112 4c7 0 11 8 11 8a18.5 18.5 0 01-2.16 3.19m-6.72-1.07a3 3 0 11-4.24-4.24"/>
                    <line x1="1" y1="1" x2="23" y2="23"/>
                  </svg>
                </button>
              </div>
              <#if messagesPerField.existsError('password-confirm')>
                <span class="cl-field-error" aria-live="polite">
                  ${kcSanitize(messagesPerField.get('password-confirm'))?no_esc}
                </span>
              </#if>
            </div>
          </#if>

        </div>

        <button tabindex="7" class="cl-primary-btn" type="submit">
          ${msg("doRegister")}
        </button>

      </form>
    </div>
  </#if>

  <#if section == "info">
    <div class="cl-footer">
      <span>${msg("haveAccount")}</span>
      <a tabindex="8" class="cl-footer-link" href="${url.loginUrl}">${msg("doLogIn")}</a>
    </div>
  </#if>

</@layout.registrationLayout>

<script>
  function togglePassword(inputId, eyeId, eyeOffId) {
    const input = document.getElementById(inputId);
    const eyeIcon = document.getElementById(eyeId);
    const eyeOffIcon = document.getElementById(eyeOffId);
    if (!input || !eyeIcon || !eyeOffIcon) return;
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
