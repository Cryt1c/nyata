<?php

namespace App\Controller;

use Symfony\Bundle\FrameworkBundle\Controller\AbstractController;
use Symfony\Component\HttpFoundation\JsonResponse;
use Symfony\Component\Routing\Annotation\Route;

#[Route('/api', name: 'api_')]
class TodoController extends AbstractController
{
    #[Route('/todos', name: 'todos', methods: ['GET'])]
    public function index(): JsonResponse
    {
        return $this->json([
            'id' => '1',
            'name' => 'Do dishes',
            'completed' => false,
        ]);
    }
}
