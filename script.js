const lawData = [
  { id: 'КЗ-101', name: 'Мелкое хулиганство / нарушение порядка', base: 'Предупреждение или короткое заключение (до 5 минут)' },
  { id: 'КЗ-202', name: 'Кража / незаконное присвоение имущества', base: 'Заключение 8–12 минут, конфискация' },
  { id: 'КЗ-303', name: 'Нападение без летального исхода', base: 'Заключение 10–15 минут, наблюдение у меда' },
  { id: 'КЗ-404', name: 'Саботаж / вред станции', base: 'Заключение 15+ минут, эскалация до пермы по контексту' },
  { id: 'КЗ-505', name: 'Тяжкое нарушение / угроза экипажу', base: 'Максимально строгое наказание по КЗ, доклад командованию' }
];

const wardenTips = [
  'Перед выдачей срока фиксируйте: статью, доказательства, свидетелей, время.',
  'Если доказательства спорные — сначала изоляция и дополнительный допрос, потом срок.',
  'Всегда уведомляйте ГСБ/капитана при серьёзных статьях или рецидиве.',
  'Сроки округляйте в меньшую сторону при явном сотрудничестве нарушителя.',
  'Ведите краткий журнал задержаний для быстрого отчёта и защиты от апелляций.'
];

const srpData = {
  'СБ': {
    'Офицер СБ': ['Нарушение процедуры задержания', 'Отсутствие отчёта после задержания', 'Превышение силы без обоснования'],
    'Варден': ['Неверная квалификация статьи', 'Нарушение порядка выдачи снаряжения', 'Неполная документация заключённого']
  },
  'Мед': {
    'Врач': ['Отказ от помощи без причины', 'Нарушение стерильных процедур', 'Неверная документация пациента'],
    'Главврач': ['Неверное распределение медперсонала', 'Отсутствие контроля критических случаев']
  },
  'Инженерия': {
    'Инженер': ['Игнорирование приказов по безопасности', 'Нарушение протокола работы с атмосферой'],
    'Старший инженер': ['Недостаточный контроль опасных работ', 'Неверное распределение доступа']
  },
  'Командование': {
    'Капитан': ['Выход за рамки СРП без протокола', 'Нарушение цепочки командования'],
    'Глава персонала': ['Ошибки в назначении доступа', 'Нарушение регламента по кадрам']
  }
};

const paperTypes = {
  'service-note': 'СЛУЖЕБНАЯ ЗАПИСКА',
  'incident-report': 'РАПОРТ О ПРОИСШЕСТВИИ',
  'ck-request': 'ЗАПРОС В ЦЕНТРАЛЬНОЕ КОМАНДОВАНИЕ',
  explanatory: 'ОБЪЯСНИТЕЛЬНАЯ'
};

const tabs = document.querySelectorAll('.tab-btn');
const panels = document.querySelectorAll('.tab-panel');
tabs.forEach(btn => btn.addEventListener('click', () => {
  tabs.forEach(b => b.classList.remove('active'));
  panels.forEach(p => p.classList.remove('active'));
  btn.classList.add('active');
  document.getElementById(btn.dataset.tab).classList.add('active');
}));

const lawSelect = document.getElementById('lawSelect');
lawData.forEach(l => {
  const o = document.createElement('option');
  o.value = l.id;
  o.textContent = `${l.id} — ${l.name}`;
  lawSelect.appendChild(o);
});

const tipsList = document.getElementById('wardenTips');
wardenTips.forEach(t => {
  const li = document.createElement('li');
  li.textContent = t;
  tipsList.appendChild(li);
});

document.getElementById('calcPunishment').addEventListener('click', () => {
  const selected = lawData.find(l => l.id === lawSelect.value);
  const aggr = document.getElementById('aggravation').value;
  const mitig = document.getElementById('mitigation').value;

  let recommendation = selected.base;
  let note = 'Базовая рекомендация.';

  if (aggr === 'minor') {
    note = 'Есть отягчающие обстоятельства: можно усилить срок на 15-25%.';
  }
  if (aggr === 'major') {
    note = 'Серьёзные отягчающие: усиление на 30-50%, уведомление командования.';
  }
  if (mitig === 'minor') {
    note += ' Смягчение: допускается понижение на 10-20%.';
  }
  if (mitig === 'major') {
    note += ' Критичное смягчение: рассмотрите альтернативу заключению.';
  }

  document.getElementById('punishmentResult').textContent =
    `Нарушение: ${selected.name}\nРекомендация: ${recommendation}\nКомментарий: ${note}`;
});

const deptSelect = document.getElementById('deptSelect');
const roleSelect = document.getElementById('roleSelect');
const srpViolation = document.getElementById('srpViolation');

Object.keys(srpData).forEach(dep => {
  const o = document.createElement('option');
  o.value = dep;
  o.textContent = dep;
  deptSelect.appendChild(o);
});

function syncRoles() {
  const dep = deptSelect.value;
  roleSelect.innerHTML = '';
  Object.keys(srpData[dep]).forEach(role => {
    const o = document.createElement('option');
    o.value = role;
    o.textContent = role;
    roleSelect.appendChild(o);
  });
  syncViolations();
}

function syncViolations() {
  const dep = deptSelect.value;
  const role = roleSelect.value;
  srpViolation.innerHTML = '';
  srpData[dep][role].forEach(v => {
    const o = document.createElement('option');
    o.value = v;
    o.textContent = v;
    srpViolation.appendChild(o);
  });
}

deptSelect.addEventListener('change', syncRoles);
roleSelect.addEventListener('change', syncViolations);
syncRoles();

document.getElementById('buildCkReport').addEventListener('click', () => {
  const now = new Date();
  const dateStr = `${String(now.getDate()).padStart(2, '0')}.${String(now.getMonth() + 1).padStart(2, '0')}.${now.getFullYear()}`;
  const text = `РАПОРТ О НАРУШЕНИИ СРП (ФОРМА CORVAX)\n`
    + `Кому: Центральное Командование\n`
    + `От кого: АВД/ПНТ\n`
    + `Дата: ${dateStr}\n`
    + `Тема: Проверка соблюдения СРП\n`
    + `Основание: Регламент СРП Corvax\n\n`
    + `1. Отдел: ${deptSelect.value}\n`
    + `2. Должность: ${roleSelect.value}\n`
    + `3. Выявленное нарушение: ${srpViolation.value}\n`
    + `4. Серьёзность: ${document.getElementById('severity').value}\n`
    + `5. Фактическая часть: ${document.getElementById('facts').value || 'не указано'}\n\n`
    + `Приложения: [заполнить при наличии]\n`
    + `Подпись: __________________`;

  document.getElementById('ckReport').textContent = text;
});

function splitList(value, fallback) {
  const parsed = value.split('\n').map(s => s.trim()).filter(Boolean);
  return parsed.length ? parsed : [fallback];
}

document.getElementById('buildInvestigation').addEventListener('click', () => {
  const caseId = document.getElementById('caseId').value || 'CASE-UNKNOWN';
  const subject = document.getElementById('subject').value || 'не указан';
  const timeline = splitList(document.getElementById('timeline').value, 'Нет данных').map((s, i) => `${i + 1}. ${s}`).join('\n');
  const evidence = splitList(document.getElementById('evidence').value, 'Нет данных').map(s => `- ${s}`).join('\n');

  document.getElementById('investigationOut').textContent =
`КАРТОЧКА РАССЛЕДОВАНИЯ АВД/ПНТ\nНомер дела: ${caseId}\nОбъект проверки: ${subject}\n\nI. Хронология событий:\n${timeline}\n\nII. Доказательства:\n${evidence}\n\nIII. Предварительный вывод:\n[заполнить после завершения проверки]\n\nIV. Рекомендация в ЦК:\n[служебная мера / доп.проверка / дисциплинарный вывод]`;
});

function getPaperFormValues() {
  return {
    type: document.getElementById('paperType').value,
    to: document.getElementById('paperTo').value.trim(),
    from: document.getElementById('paperFrom').value.trim(),
    date: document.getElementById('paperDate').value.trim(),
    subject: document.getElementById('paperSubject').value.trim(),
    basis: document.getElementById('paperBasis').value.trim(),
    body: document.getElementById('paperBody').value.trim(),
    attachments: document.getElementById('paperAttachments').value.trim()
  };
}

function validatePaper(values) {
  const missing = [];
  if (!values.to) missing.push('Кому');
  if (!values.from) missing.push('От кого');
  if (!values.date) missing.push('Дата / смена');
  if (!values.subject) missing.push('Тема');
  if (!values.basis) missing.push('Основание');
  if (!values.body) missing.push('Содержание');
  return missing;
}

function generatePaperByForm(values) {
  return `${paperTypes[values.type]}\n`
    + `Кому: ${values.to}\n`
    + `От кого: ${values.from}\n`
    + `Дата / смена: ${values.date}\n`
    + `Тема: ${values.subject}\n`
    + `Основание: ${values.basis}\n\n`
    + `Содержание:\n${values.body}\n\n`
    + `Приложения: ${values.attachments || 'отсутствуют'}\n\n`
    + `Подпись: __________________\n`
    + `Расшифровка подписи: __________________`;
}

document.getElementById('validatePaperwork').addEventListener('click', () => {
  const values = getPaperFormValues();
  const missing = validatePaper(values);
  const status = document.getElementById('paperValidation');
  if (missing.length) {
    status.textContent = `Форма заполнена НЕ полностью. Обязательные поля: ${missing.join(', ')}.`;
    status.className = 'validation bad';
    return;
  }
  status.textContent = 'Форма заполнена корректно: обязательные реквизиты присутствуют.';
  status.className = 'validation good';
});

document.getElementById('generatePaperwork').addEventListener('click', () => {
  const values = getPaperFormValues();
  const missing = validatePaper(values);
  if (missing.length) {
    document.getElementById('paperValidation').textContent = `Невозможно сформировать документ: заполните поля ${missing.join(', ')}.`;
    document.getElementById('paperValidation').className = 'validation bad';
    return;
  }
  document.getElementById('paperValidation').textContent = 'Документ сформирован строго по форме с обязательными реквизитами.';
  document.getElementById('paperValidation').className = 'validation good';
  document.getElementById('paperworkOut').textContent = generatePaperByForm(values);
});
