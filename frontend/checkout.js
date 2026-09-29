// checkout.js
(function () {
  const emailInput    = document.getElementById('email');
  const phoneInput    = document.getElementById('phone');
  const addressInput  = document.getElementById('shippingAddress');
  const checkoutForm  = document.getElementById('checkoutForm');
  const submitBtn     = document.getElementById('submitBtn');
  const successMsg    = document.getElementById('successMsg');
  const welcomeBanner = document.getElementById('welcomeBanner');
  const recognizeStatus = document.getElementById('recognizeStatus');

  // Modal elements
  const otpModal   = document.getElementById('otpModal');
  const modalEmail = document.getElementById('modalEmail');
  const otpInput   = document.getElementById('otpInput');
  const otpError   = document.getElementById('otpError');
  const verifyBtn  = document.getElementById('verifyBtn');
  const skipBtn    = document.getElementById('skipBtn');

  const emailRe = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  const phoneRe = /^\+?[\d\s\-().]{7,20}$/;

  let loggedInUser = null;
  let recognizeTimer = null;
  let lastRecognizedEmail = '';

  // ===== Validation helpers =====
  function setError(inputEl, errorId, msg) {
    document.getElementById(errorId).textContent = msg;
    if (msg) inputEl.classList.add('invalid');
    else inputEl.classList.remove('invalid');
  }

  function clearFormErrors() {
    setError(emailInput,   'emailError',   '');
    setError(phoneInput,   'phoneError',   '');
    setError(addressInput, 'addressError', '');
  }

  // ===== Email recognition (debounced) =====
  emailInput.addEventListener('input', () => {
    clearTimeout(recognizeTimer);
    const email = emailInput.value.trim();

    if (!emailRe.test(email)) {
      recognizeStatus.textContent = '';
      return;
    }

    recognizeStatus.textContent = 'Checking…';
    recognizeTimer = setTimeout(() => runRecognition(email), 600);
  });

  async function runRecognition(email) {
    if (email === lastRecognizedEmail) return;
    lastRecognizedEmail = email;
    try {
      const { registered } = await api.recognize(email);
      recognizeStatus.textContent = '';
      if (registered && !loggedInUser) {
        openModal(email);
      }
    } catch {
      recognizeStatus.textContent = '';
    }
  }

  // ===== Modal =====
  function openModal(email) {
    modalEmail.textContent = email;
    otpInput.value = '';
    otpError.textContent = '';
    otpModal.classList.remove('hidden');
    setTimeout(() => otpInput.focus(), 50);
  }

  function closeModal() {
    otpModal.classList.add('hidden');
  }

  skipBtn.addEventListener('click', closeModal);

  // Close modal on overlay click
  otpModal.addEventListener('click', (e) => {
    if (e.target === otpModal) closeModal();
  });

  // Only allow digits in OTP input
  otpInput.addEventListener('input', () => {
    otpInput.value = otpInput.value.replace(/\D/g, '').slice(0, 6);
    otpError.textContent = '';
  });

  verifyBtn.addEventListener('click', handleVerify);
  otpInput.addEventListener('keydown', (e) => { if (e.key === 'Enter') handleVerify(); });

  function setVerifyLoading(loading) {
    verifyBtn.disabled = loading;
    verifyBtn.querySelector('.btn-text').textContent = loading ? 'Verifying…' : 'Verify & Login';
    verifyBtn.querySelector('.spinner').classList.toggle('hidden', !loading);
  }

  async function handleVerify() {
    const code = otpInput.value.trim();
    if (code.length !== 6) {
      otpError.textContent = 'Please enter a 6-digit code';
      return;
    }
    if (!/^\d{6}$/.test(code)) {
      otpError.textContent = 'Code must contain digits only';
      return;
    }

    setVerifyLoading(true);
    try {
      const data = await api.verify(emailInput.value.trim(), code);
      loggedInUser = data;
      closeModal();
      showWelcome(data.firstName, data.lastName);
    } catch (err) {
      otpError.textContent = err.message || 'Invalid code. Please try again.';
    } finally {
      setVerifyLoading(false);
    }
  }

  function showWelcome(firstName, lastName) {
    welcomeBanner.textContent = `Welcome, ${firstName} ${lastName} 👋`;
    welcomeBanner.classList.remove('hidden');
    welcomeBanner.scrollIntoView({ behavior: 'smooth', block: 'nearest' });
  }

  // ===== Checkout form submission =====
  function validateCheckout() {
    let valid = true;
    const email   = emailInput.value.trim();
    const phone   = phoneInput.value.trim();
    const address = addressInput.value.trim();

    if (!email) {
      setError(emailInput, 'emailError', 'Email is required');
      valid = false;
    } else if (!emailRe.test(email)) {
      setError(emailInput, 'emailError', 'Please enter a valid email address');
      valid = false;
    } else {
      setError(emailInput, 'emailError', '');
    }

    if (!phone) {
      setError(phoneInput, 'phoneError', 'Phone number is required');
      valid = false;
    } else if (!phoneRe.test(phone)) {
      setError(phoneInput, 'phoneError', 'Please enter a valid phone number');
      valid = false;
    } else {
      setError(phoneInput, 'phoneError', '');
    }

    if (!address) {
      setError(addressInput, 'addressError', 'Shipping address is required');
      valid = false;
    } else {
      setError(addressInput, 'addressError', '');
    }

    return valid;
  }

  function setSubmitLoading(loading) {
    submitBtn.disabled = loading;
    submitBtn.querySelector('.btn-text').textContent = loading ? 'Submitting…' : 'Place Order';
    submitBtn.querySelector('.spinner').classList.toggle('hidden', !loading);
  }

  checkoutForm.addEventListener('submit', async (e) => {
    e.preventDefault();
    clearFormErrors();
    successMsg.classList.add('hidden');

    if (!validateCheckout()) return;

    setSubmitLoading(true);
    try {
      await api.checkout(
        emailInput.value.trim(),
        phoneInput.value.trim(),
        addressInput.value.trim()
      );
      successMsg.textContent = '✅ Checkout submitted successfully! Your order is confirmed.';
      successMsg.classList.remove('hidden');
      checkoutForm.reset();
      loggedInUser = null;
      welcomeBanner.classList.add('hidden');
      lastRecognizedEmail = '';
      successMsg.scrollIntoView({ behavior: 'smooth', block: 'nearest' });
    } catch (err) {
      successMsg.textContent = '';
      setError(emailInput, 'emailError', err.message || 'Checkout failed. Please try again.');
    } finally {
      setSubmitLoading(false);
    }
  });
})();
