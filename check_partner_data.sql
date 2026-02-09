-- Check Ahmad Pratama's data
SELECT 'Partner Info:' as info;
SELECT id, referral_code, total_referrals, total_commission, available_balance, paid_amount 
FROM referral_partners 
WHERE referral_code = 'AHMAD2024';

SELECT '\nCompanies Referred:' as info;
SELECT COUNT(*) as count 
FROM partner_referrals pr
JOIN referral_partners rp ON pr.partner_id = rp.id
WHERE rp.referral_code = 'AHMAD2024';

SELECT '\nCommissions Earned:' as info;
SELECT COUNT(*) as count, COALESCE(SUM(commission_amount), 0) as total
FROM partner_commissions pc
JOIN referral_partners rp ON pc.partner_id = rp.id
WHERE rp.referral_code = 'AHMAD2024';
