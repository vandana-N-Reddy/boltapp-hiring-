// register.js
(function () {
  const form        = document.getElementById('registerForm');
  const emailInput  = document.getElementById('email');
  const firstInput  = document.getElementById('firstName');
  const lastInput   = document.getElementById('lastName');
  const registerBtn = document.getElementById('registerBtn');
  const successPanel = document.getElementById('successPanel');
  const otpCode     = document.getElementById('otpCode');
  const errorMsg    = document.getElementById('errorMsg');

  const emailRe = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

  function setError(inputEl, errorId, msg) {
    const el = document.getElementById(errorId);
    el.textContent = msg;
    if (msg) inputEl.classList.add('invalid');
    else inputEl.classList.remove('invalid');
  }

  function clearErrors() {
    ['emailError','firstNameError','lastNameError'].forEach(id => {
      document.getElementById(id).textContent = '';
    });
    [emailInput, firstInput, lastInput].forEach(el => el.classList.remove('invalid'));
    errorMsg.classList.add('hidden');
    errorMsg.textContent = '';
  }

  function validate() {
    let valid = true;
    const email = emailInput.value.trim();
    const first = firstInput.value.trim();
    const last  = lastInput.value.trim();

    if (!email) {
      setError(emailInput, 'emailError', 'Email is required');
      valid = false;
    } else if (!emailRe.test(email)) {
      setError(emailInput, 'emailError', 'Please enter a valid email address');
      valid = false;
    } else {
      setError(emailInput, 'emailError', '');
    }

    if (!first) {
      setError(firstInput, 'firstNameError', 'First name is required');
      valid = false;
    } else {
      setError(firstInput, 'firstNameError', '');
    }

    if (!last) {
      setError(lastInput, 'lastNameError', 'Last name is required');
      valid = false;
    } else {
      setError(lastInput, 'lastNameError', '');
    }

    return valid;
  }

  function setLoading(loading) {
    registerBtn.disabled = loading;
    registerBtn.querySelector('.btn-text').textContent = loading ? 'Registering…' : 'Register';
    registerBtn.querySelector('.spinner').classList.toggle('hidden', !loading);
  }

  form.addEventListener('submit', async (e) => {
    e.preventDefault();
    clearErrors();
    if (!validate()) return;

    setLoading(true);
    try {
      const data = await api.register(
        emailInput.value.trim(),
        firstInput.value.trim(),
        lastInput.value.trim()
      );
      form.classList.add('hidden');
      otpCode.textContent = data.code;
      successPanel.classList.remove('hidden');
    } catch (err) {
      errorMsg.textContent = err.message;
      errorMsg.classList.remove('hidden');
      if (err.message.toLowerCase().includes('already registered')) {
        setError(emailInput, 'emailError', 'Email is already registered');
      }
    } finally {
      setLoading(false);
    }
  });
})();
