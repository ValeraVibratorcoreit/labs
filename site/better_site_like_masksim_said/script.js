const fortunes = [
  "Дорога начнётся с телефонного звонка, который вы чуть не пропустите.",
  "Тот, кого вы давно не видели, окажется рядом раньше, чем вы думаете.",
  "Решение, которое вы откладывали, примет себя само — и правильно.",
  "Ждите новости в среду. Возможно, не в эту.",
  "Кто-то запомнит сегодняшний разговор с вами надолго.",
  "Дно чашки чистое — это значит, неделя пройдёт без сюрпризов.",
  "Впереди — мелкая удача, но именно та, что нужна сейчас.",
  "Не открывайте письмо сразу. Дайте ему полежать день.",
  "Вы близки к ответу на вопрос, который ещё не задали вслух.",
  "Гуща сложилась в форму двери. Толкните, а не тяните."
];

const btn = document.getElementById('fortune-btn');
const out = document.getElementById('fortune-text');

btn.addEventListener('click', () => {
  btn.disabled = true;
  const pick = fortunes[Math.floor(Math.random() * fortunes.length)];
  
  out.classList.remove('show');
  
  window.setTimeout(() => {
    out.classList.remove('placeholder');
    out.textContent = pick;
    out.classList.add('show');
    btn.disabled = false;
  }, 180);
});