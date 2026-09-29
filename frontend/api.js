// api.js — centralised API calls
// Reads API_BASE from window.API_BASE (set via env config) or defaults to localhost

const API_BASE = window.API_BASE || 'http://localhost:8080';

const api = {
  async register(email, firstName, lastName) {
    const res = await fetch(`${API_BASE}/api/auth/register`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, firstName, lastName }),
    });
    const data = await res.json();
    if (!res.ok) throw new Error(data.error || 'Registration failed');
    return data;
  },

  async recognize(email) {
    const res = await fetch(`${API_BASE}/api/users/recognize?email=${encodeURIComponent(email)}`);
    const data = await res.json();
    if (!res.ok) throw new Error(data.error || 'Recognition failed');
    return data; // { registered: true/false }
  },

  async verify(email, code) {
    const res = await fetch(`${API_BASE}/api/auth/verify`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, code }),
    });
    const data = await res.json();
    if (!res.ok) throw new Error(data.error || 'Verification failed');
    return data;
  },

  async checkout(email, phone, shippingAddress) {
    const res = await fetch(`${API_BASE}/api/checkout`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, phone, shippingAddress }),
    });
    const data = await res.json();
    if (!res.ok) throw new Error(data.error || 'Checkout failed');
    return data;
  },
};
