<html>
<head>
  <title>Contact Me</title>
  <style>
    body {
      margin: 0;
      padding: 0;
      min-height: 100vh;
      background-color: #f8f8f8;
      display: flex;
      flex-direction: column;
    }

    main {
      flex: 1;
    }

    .contact {
      padding: 80px 0;
    }

    .container {
      max-width: 800px;
      margin: 0 auto;
      padding: 0 20px;
    }

    h1 {
      font-size: 36px;
      font-weight: bold;
      text-align: center;
      margin-bottom: 20px;
      color: #333;
    }

    p {
      font-size: 18px;
      text-align: center;
      margin-bottom: 40px;
      color: #666;
    }

    .contact-info {
      display: flex;
      justify-content: center;
    }

    .info-item {
      display: flex;
      align-items: center;
      padding: 30px;
      background-color: #fff;
      border-radius: 8px;
      box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
      transition: transform 0.3s ease;
    }

    .info-item:hover {
      transform: translateY(-5px);
    }

    .icon {
      font-size: 36px;
      margin-right: 20px;
      color: #007bff;
    }

    .details h3 {
      font-size: 24px;
      margin-bottom: 10px;
      color: #333;
    }

    .details p {
      font-size: 18px;
      margin: 0;
      color: #666;
      text-align: left;
    }
  </style>
</head>
<body>
  <main>
    <section class="contact">
      <div class="container">
        <h1>Contact Me</h1>
        <p>Feel free to reach out to me using the email address below.</p>

        <div class="contact-info">
          <div class="info-item">
            <div class="icon">
              <i class="fas fa-envelope"></i>
            </div>
            <div class="details">
              <h3>Email</h3>
              <p>yourname@example.com</p>
            </div>
          </div>
        </div>
      </div>
    </section>
  </main>
</body>
</html>
